package fs

import (
	"context"
	"io"
	"os"
	"path"
	"path/filepath"

	"github.com/qingstor/go-mime"

	"github.com/aos-dev/go-storage/v2/pkg/iowrap"
	"github.com/aos-dev/go-storage/v2/types"
	"github.com/aos-dev/go-storage/v2/types/info"
)

func (s *Storage) delete(ctx context.Context, path string, opt *pairStorageDelete) (err error) {
	rp := s.getAbsPath(path)

	err = s.osRemove(rp)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) listDir(ctx context.Context, dir string, opt *pairStorageListDir) (err error) {
	// Always keep service original name as rp.
	rp := s.getAbsPath(dir)
	// Then convert the dir to slash separator.
	dir = filepath.ToSlash(dir)

	fi, err := s.ioutilReadDir(rp)
	if err != nil {
		return err
	}

	for _, v := range fi {
		// if v is a link, and client not follow link, skip it
		if v.Mode()&os.ModeSymlink != 0 && !opt.EnableLinkFollow {
			continue
		}

		target, err := checkLink(v, rp)
		if err != nil {
			return err
		}

		o := &types.Object{
			// Always keep service original name as ID.
			ID: filepath.Join(rp, v.Name()),
			// Object's name should always be separated by slash (/)
			Name:       path.Join(dir, v.Name()),
			Size:       target.Size(),
			UpdatedAt:  target.ModTime(),
			ObjectMeta: info.NewObjectMeta(),
		}

		if target.IsDir() {
			o.Type = types.ObjectTypeDir
			if opt.HasDirFunc {
				opt.DirFunc(o)
			}
			continue
		}

		if v := mime.DetectFilePath(target.Name()); v != "" {
			o.SetContentType(v)
		}

		o.Type = types.ObjectTypeFile
		if opt.HasFileFunc {
			opt.FileFunc(o)
		}
	}
	return
}

func (s *Storage) metadata(ctx context.Context, opt *pairStorageMetadata) (meta info.StorageMeta, err error) {
	meta = info.NewStorageMeta()
	meta.WorkDir = s.workDir
	return meta, nil
}

func (s *Storage) read(ctx context.Context, path string, opt *pairStorageRead) (rc io.ReadCloser, err error) {
	// If path is "-", return stdin directly.
	if path == "-" {
		f := os.Stdin
		if opt.HasSize {
			return iowrap.LimitReadCloser(f, opt.Size), nil
		}
		return f, nil
	}

	rp := s.getAbsPath(path)

	f, err := s.osOpen(rp)
	if err != nil {
		return nil, err
	}
	if opt.HasOffset {
		_, err = f.Seek(opt.Offset, 0)
		if err != nil {
			return nil, err
		}
	}

	rc = f
	if opt.HasSize {
		rc = iowrap.LimitReadCloser(rc, opt.Size)
	}
	if opt.HasReadCallbackFunc {
		rc = iowrap.CallbackReadCloser(rc, opt.ReadCallbackFunc)
	}
	return rc, nil
}

func (s *Storage) stat(ctx context.Context, path string, opt *pairStorageStat) (o *types.Object, err error) {
	if path == "-" {
		return &types.Object{
			ID:         "-",
			Name:       "-",
			Type:       types.ObjectTypeStream,
			Size:       0,
			ObjectMeta: info.NewObjectMeta(),
		}, nil
	}

	rp := s.getAbsPath(path)

	fi, err := s.osStat(rp)
	if err != nil {
		return nil, err
	}

	o = &types.Object{
		ID:         rp,
		Name:       path,
		Size:       fi.Size(),
		UpdatedAt:  fi.ModTime(),
		ObjectMeta: info.NewObjectMeta(),
	}

	if fi.IsDir() {
		o.Type = types.ObjectTypeDir
		return
	}
	if fi.Mode().IsRegular() {
		if v := mime.DetectFilePath(path); v != "" {
			o.SetContentType(v)
		}

		o.Type = types.ObjectTypeFile
		return
	}
	if fi.Mode()&StreamModeType != 0 {
		o.Type = types.ObjectTypeStream
		return
	}

	o.Type = types.ObjectTypeInvalid
	return o, nil
}

func (s *Storage) write(ctx context.Context, path string, r io.Reader, opt *pairStorageWrite) (err error) {
	var f io.WriteCloser
	// If path is "-", use stdout directly.
	if path == "-" {
		f = os.Stdout
	} else {
		// Create dir for path.
		err = s.createDir(path)
		if err != nil {
			return err
		}

		rp := s.getAbsPath(path)

		f, err = s.osCreate(rp)
		if err != nil {
			return err
		}
	}

	if opt.HasReadCallbackFunc {
		r = iowrap.CallbackReader(r, opt.ReadCallbackFunc)
	}

	if opt.HasSize {
		_, err = s.ioCopyN(f, r, opt.Size)
	} else {
		_, err = s.ioCopyBuffer(f, r, make([]byte, 1024*1024))
	}
	if err != nil {
		return err
	}
	return
}

func (s *Storage) copy(ctx context.Context, src string, dst string, opt *pairStorageCopy) (err error) {
	rs := s.getAbsPath(src)
	rd := s.getAbsPath(dst)

	// Create dir for dst.
	err = s.createDir(dst)
	if err != nil {
		return err
	}

	srcFile, err := s.osOpen(rs)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := s.osCreate(rd)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = s.ioCopyBuffer(dstFile, srcFile, make([]byte, 1024*1024))
	if err != nil {
		return err
	}
	return
}
func (s *Storage) move(ctx context.Context, src string, dst string, opt *pairStorageMove) (err error) {

	rs := s.getAbsPath(src)
	rd := s.getAbsPath(dst)

	// Create dir for dst path.
	err = s.createDir(dst)
	if err != nil {
		return err
	}

	err = s.osRename(rs, rd)
	if err != nil {
		return err
	}
	return
}

func checkLink(v os.FileInfo, dir string) (os.FileInfo, error) {
	// if v is not link, return directly
	if v.Mode()&os.ModeSymlink == 0 {
		return v, nil
	}

	// otherwise, follow the link to get the target
	tarPath, err := filepath.EvalSymlinks(filepath.Join(dir, v.Name()))
	if err != nil {
		return nil, err
	}
	return os.Stat(tarPath)
}
