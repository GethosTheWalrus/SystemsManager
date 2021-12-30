package controllers

import (
	"ThePooReview/models"
	"encoding/json"
	"net/http"
)

func GetSystems(w http.ResponseWriter, r *http.Request) {
	systems := []models.HostSystem{}
	Db.
		Preload("Cpu").
		Preload("Memory").
		Preload("Network").
		Preload("General").
		Find(&systems)

	json.NewEncoder(w).Encode(systems)

}
