package nursery

import (
   "archive/tar"
   "archive/zip"
   "github.com/xi2/xz"
   "io"
   "os"
   "path/filepath"
)

func Xz(in, dir string, level int) error {
   file, err := os.Open(in)
   if err != nil {
      return err
   }
   defer file.Close()
   read, err := xz.NewReader(file, 0)
   if err != nil {
      return err
   }
   return Tar(read, dir, level)
}

func strip(in string, n int) string {
   for i, char := range in {
      if n <= 0 {
         return in[i:]
      }
      if char == '/' {
         n--
      }
   }
   return ""
}

// Need for Zstandard
func Tar(in io.Reader, dir string, level int) error {
   read := tar.NewReader(in)
   for {
      head, err := read.Next()
      if err == io.EOF {
         break
      } else if err != nil {
         return err
      }
      head.Name = strip(head.Name, level)
      if head.Name != "" {
         head.Linkname = strip(head.Linkname, level)
         err := create(head, read, dir)
         if err != nil {
            return err
         }
      }
   }
   return nil
}

func create(head *tar.Header, in io.Reader, out string) error {
   head.Name = filepath.Join(out, head.Name)
   err := os.MkdirAll(filepath.Dir(head.Name), os.ModePerm)
   if err != nil {
      return err
   }
   switch head.Typeflag {
   case tar.TypeReg:
      file, err := os.Create(head.Name)
      if err != nil {
         return err
      }
      defer file.Close()
      if _, err := file.ReadFrom(in); err != nil {
         return err
      }
   case tar.TypeLink:
      _, err := os.Stat(head.Name)
      if err != nil {
         err := os.Link(filepath.Join(out, head.Linkname), head.Name)
         if err != nil {
            return err
         }
      }
   }
   return nil
}

func Zip(in, dir string, level int) error {
   read, err := zip.OpenReader(in)
   if err != nil {
      return err
   }
   defer read.Close()
   for _, head := range read.File {
      if head.Mode().IsDir() {
         continue
      }
      in, err := head.Open()
      if err != nil {
         return err
      }
      defer in.Close()
      head.Name = filepath.Join(dir, strip(head.Name, level))
      if err := os.MkdirAll(filepath.Dir(head.Name), os.ModePerm); err != nil {
         return err
      }
      out, err := os.Create(head.Name)
      if err != nil {
         return err
      }
      defer out.Close()
      if _, err := out.ReadFrom(in); err != nil {
         return err
      }
   }
   return nil
}
