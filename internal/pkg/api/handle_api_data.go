package api

import (
	"encoding/json"
	"net/http"

	"github.com/dgraph-io/badger/v3"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/zltl/nydus-auth/internal/pkg/m"
	"github.com/zltl/nydus-auth/pkg/id"
)

func (s *State) handleApiGetData(c *gin.Context) {
	uid := c.MustGet("uid").(id.ID)

	// s.kvDB.Get
	err := s.kvDB.View(func(txn *badger.Txn) error {
		item, err := txn.Get(uid.Bytes())
		if err != nil {
			logrus.Error(err)
			return err
		}
		err = item.Value(func(val []byte) error {
			var ds m.StoreData
			err := json.Unmarshal(val, &ds)
			if err != nil {
				logrus.Error(err)
				return err
			}
			c.JSON(http.StatusOK, gin.H{
				"uid":       uid,
				"timestamp": ds.Ts,
				"item":      ds.Data,
				"from":      ds.From,
			})
			return nil
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":             "internal_error",
				"error_description": "could not get data",
			})
			return err
		}
		return nil
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":             "internal_error",
			"error_description": "could not get data",
		})
	}
}

func (s *State) handleApiPostData(c *gin.Context) {
	uid := c.MustGet("uid").(id.ID)

	var ds m.StoreData
	dsBytes, err := c.GetRawData()
	if err != nil {
		logrus.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":             "invalid_request",
			"error_description": "could not parse request body",
		})
		return
	}
	err = json.Unmarshal(dsBytes, &ds)
	if err != nil {
		logrus.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":             "invalid_request",
			"error_description": "could not parse request body",
		})
		return
	}

	err = s.kvDB.Update(func(txn *badger.Txn) error {
		e := badger.NewEntry(uid.Bytes(), dsBytes)
		err := txn.SetEntry(e)
		if err != nil {
			logrus.Error(err)
			return err
		}
		return nil
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":             "internal_error",
			"error_description": "could not set data",
		})
		return
	}
}
