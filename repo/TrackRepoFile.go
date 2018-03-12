package repo

import (
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/x1m3/geotracker/entity"
	"io"
	"os"
	"sort"
	"sync"
	"time"
)

type TrackRepoFile struct {
	sync.Mutex
	fp *os.File
}

func NewTracRepoFile(fp *os.File) *TrackRepoFile {
	return &TrackRepoFile{fp: fp}
}

func (r *TrackRepoFile) Store(track *entity.Track) error {
	r.Lock()
	defer r.Unlock()

	dto := NewtrackFileDTO(track)
	err := binary.Write(r.fp, binary.BigEndian, dto)
	return err
}

func (r *TrackRepoFile) GetTracksByDriverAsc(driverID int64) (entity.TrackList, error) {
	var err error = nil
	var dto *TrackFileDTO

	r.Lock()
	defer func() {
		r.fp.Seek(0, io.SeekEnd)
		r.Unlock()
	}()
	tracks := make([]*entity.Track, 0, 0)
	r.fp.Seek(0, io.SeekStart)

	for {
		dto, err = r.readRecord()
		switch err {
		case nil:
			if tracks, err = r.filterByDriverID(tracks, dto, driverID); err != nil {
				return nil, err
			}
		case io.EOF:
			return tracks, nil
		default:
			return nil, err
		}
	}
}

func (r *TrackRepoFile) filterByDriverID(tracks entity.TrackList, dto *TrackFileDTO, driverID int64) (entity.TrackList, error) {
	if dto.DriverID == driverID {
		c, err := entity.NewCoordinate(dto.Lat, dto.Lon)
		if err != nil {
			return nil, errors.New("reading wrong coordinate")
		}
		tracks = append(tracks, entity.NewTrack(c, driverID, time.Unix(0, dto.TimeUnix)))
	}
	sort.Sort(entity.ByDate(tracks))
	return tracks, nil
}

func (r *TrackRepoFile) readRecord() (*TrackFileDTO, error) {
	dto := &TrackFileDTO{}
	err := binary.Read(r.fp, binary.BigEndian, dto)
	return dto, err
}

type TrackFileDTO struct {
	DriverID int64
	Lat      float64
	Lon      float64
	TimeUnix int64
}

func NewtrackFileDTO(t *entity.Track) *TrackFileDTO {
	t.ReceivedOn.UnixNano()
	return &TrackFileDTO{
		DriverID: t.DriverID,
		Lat:      t.Point.Lat(),
		Lon:      t.Point.Lon(),
		TimeUnix: t.ReceivedOn.UnixNano(),
	}
}

func (dto *TrackFileDTO) serialize() []byte {
	var buff bytes.Buffer

	binary.Write(&buff, binary.BigEndian, dto)
	return buff.Bytes()
}
