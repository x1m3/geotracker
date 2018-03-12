package repo

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/x1m3/geotracker/entity"
	"log"
	"time"
)

/*
CREATE TABLE `track` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `driver_id` int(11) NOT NULL,
  `lat` double NOT NULL,
  `lon` double NOT NULL,
  `created_on` timestamp(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
  PRIMARY KEY (`id`),
  KEY `track_by_driver_idx` (`driver_id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;

*/

type TrackMysqlDTO struct {
	DriverID   int64
	Lat        float64
	Lon        float64
	Created_on time.Time
}

type TrackRepoMYSQL struct {
	db                  *sql.DB
	storeTrackPS        *sql.Stmt
	tracksByDriverAscPS *sql.Stmt
}

func NewTrackRepoMYSQL(db *sql.DB) *TrackRepoMYSQL {
	var err error
	repo := &TrackRepoMYSQL{db: db}
	repo.storeTrackPS, err = repo.db.Prepare("INSERT INTO track (driver_id, lat, lon, created_on) VALUES(?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	repo.tracksByDriverAscPS, err = repo.db.Prepare("SELECT driver_id, lat, lon, created_on FROM track WHERE driver_id=? ORDER BY created_on, id ASC")
	if err != nil {
		log.Fatal(err)
	}
	return repo
}

func (r *TrackRepoMYSQL) Store(track *entity.Track) error {
	_, err := r.storeTrackPS.Exec(track.DriverID, track.Point.Lat(), track.Point.Lon(), track.ReceivedOn)
	return err
}

func (r *TrackRepoMYSQL) GetTracksByDriverAsc(driverID int64) (entity.TrackList, error) {
	rows, err := r.tracksByDriverAscPS.Query(driverID)
	defer rows.Close()

	dto := TrackMysqlDTO{}
	tracks := make(entity.TrackList, 0)

	switch err {
	case nil:
		for rows.Next() {
			rows.Scan(&dto.DriverID, &dto.Lat, &dto.Lon, &dto.Created_on)
			c, err := entity.NewCoordinate(dto.Lat, dto.Lon)
			if err != nil {
				return nil, err
			}
			tracks = append(tracks, entity.NewTrack(c, dto.DriverID, dto.Created_on))
		}
	case sql.ErrNoRows:
		return tracks, nil
	default:
		return nil, err
	}
	return tracks, nil
}
