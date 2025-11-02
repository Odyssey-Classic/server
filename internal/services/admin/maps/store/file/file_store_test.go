package file

import (
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
)

type FileStoreSuite struct {
	suite.Suite
	dir string
	fs  *FileStore
}

func (s *FileStoreSuite) SetupTest() {
	d, err := os.MkdirTemp("", "maps-store-test-*")
	s.Require().NoError(err)
	s.dir = d
	fs, err := New(d)
	s.Require().NoError(err)
	s.fs = fs
}

func (s *FileStoreSuite) TearDownTest() {
	os.RemoveAll(s.dir)
}

func (s *FileStoreSuite) TestCreateAndGet() {
	m, err := s.fs.Create("Alpha")
	s.Require().NoError(err)
	s.NotZero(m.ID)

	got, err := s.fs.Get(m.ID)
	s.Require().NoError(err)
	s.Equal("Alpha", got.Name)
}

func (s *FileStoreSuite) TestListAndSearch() {
	s.fs.Create("Alpha")
	s.fs.Create("Beta")
	s.fs.Create("Gamma")

	all, err := s.fs.List("")
	s.Require().NoError(err)
	s.Len(all, 3)

	b, err := s.fs.List("et")
	s.Require().NoError(err)
	s.Len(b, 1)
	s.Equal("Beta", b[0].Name)
}

func (s *FileStoreSuite) TestUpdate() {
	m, _ := s.fs.Create("Alpha")
	m.Name = "Alpha Prime"
	s.Require().NoError(s.fs.Update(m))
	got, _ := s.fs.Get(m.ID)
	s.Equal("Alpha Prime", got.Name)
}

func (s *FileStoreSuite) TestDelete() {
	m, _ := s.fs.Create("Alpha")
	s.Require().NoError(s.fs.Delete(m.ID))
	_, err := s.fs.Get(m.ID)
	s.Error(err)
}

func TestFileStore(t *testing.T) {
	suite.Run(t, new(FileStoreSuite))
}
