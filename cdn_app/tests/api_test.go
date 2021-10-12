package tests

import (
	"cdn_app/fixtures"
	api "cdn_app/pkg/files/http"
	fservice "cdn_app/pkg/files/service"
	fstorage "cdn_app/pkg/files/storage"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAPIHandler(t *testing.T) {
	rsc, err := fixtures.InitTestDB()
	if err != nil {
		fmt.Printf("Can't initialize resources. Err: %v", err)
	}
	defer func() {
		err := rsc.Release()
		if err != nil {
			fmt.Printf("Got an error during resources release. %v", err)
		}
	}()

	rsc.DB.Exec(`
		INSERT INTO area (name, description) VALUES ('test.ru', 'TST');
		INSERT INTO file (name, sha, size, description) VALUES
			('file1.mp4', '123456789', 1024, 'desc file1.mp4'),
			('file2.mp4', '678912345', 4096, 'desc file2.mp4');
		INSERT INTO server (area_id, name, hostname, description) VALUES
			(1, 'srv1', 'srv1.test.ru', 'desc srv1.test.ru'),
			(1, 'srv2', 'srv2.test.ru', 'desc srv2.test.ru');
		INSERT INTO "user" (name, balance) VALUES
			('user1', 100),
			('user2', 200);
		INSERT INTO userfile (user_id, file_id) VALUES
			(1, 1),
			(1, 2),
			(2, 2);
		INSERT INTO serverfile (server_id, file_id) VALUES
			(1, 1),
			(1, 2),
			(2, 1);
	`)

	filesDb := fstorage.New(rsc.DB)
	filesCtrl := fservice.NewController(filesDb)
	apiHandler := api.APIHandler(filesCtrl)

	t.Run("Case Get user files", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/api/files/user/user1", nil)
		rw := httptest.NewRecorder()
		vars := map[string]string{
			"name": "user1",
		}
		req = mux.SetURLVars(req, vars)
		apiHandler.ServeHTTP(rw, req)
		assert.Equal(t, rw.Code, http.StatusOK)
		var respBody []string
		json.NewDecoder(rw.Body).Decode(&respBody)
		assert.Equal(t, respBody, []string{"file1.mp4", "file2.mp4"})
	})

	t.Run("Case Get server files", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/api/files/server/srv2", nil)
		rw := httptest.NewRecorder()
		vars := map[string]string{
			"name": "srv2",
		}
		req = mux.SetURLVars(req, vars)
		apiHandler.ServeHTTP(rw, req)
		assert.Equal(t, rw.Code, http.StatusOK)
		var respBody []string
		json.NewDecoder(rw.Body).Decode(&respBody)
		assert.Equal(t, respBody, []string{"file1.mp4"})
	})

	t.Run("Case Get area files", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/api/files/area/test.ru", nil)
		rw := httptest.NewRecorder()
		vars := map[string]string{
			"name": "test.ru",
		}
		req = mux.SetURLVars(req, vars)
		apiHandler.ServeHTTP(rw, req)
		assert.Equal(t, rw.Code, http.StatusOK)
		var respBody []string
		json.NewDecoder(rw.Body).Decode(&respBody)
		assert.Equal(t, respBody, []string{"file1.mp4", "file2.mp4"})
	})

}
