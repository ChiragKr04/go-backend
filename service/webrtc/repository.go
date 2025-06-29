package webrtc

import (
	"ChiragKr04/go-backend/types"
	"database/sql"
)

type WebrtcRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *WebrtcRepository {
	return &WebrtcRepository{db: db}
}

func (r *WebrtcRepository) CreateOffer(offer types.Offer) (types.Offer, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return types.Offer{}, err
	}

	offerResult, err := tx.Exec("INSERT INTO offers (offer, offerer_user_id, room_id, offer_ice_candidates) VALUES (?, ?, ?, ?)", offer.Offer, offer.OffererUserID, offer.RoomID, offer.OfferIceCandidates)
	if err != nil {
		_ = tx.Rollback()
		return types.Offer{}, err
	}

	_, err = offerResult.LastInsertId()	
	if err != nil {
		_ = tx.Rollback()
		return types.Offer{}, err
	}

	err = tx.Commit()
	if err != nil {
		return types.Offer{}, err
	}

	return types.Offer{
		Offer:               offer.Offer,
		OffererUserID:       offer.OffererUserID,
		RoomID:              offer.RoomID,
		OfferIceCandidates:  offer.OfferIceCandidates,
		AnswerIceCandidates: offer.AnswerIceCandidates,
	}, nil
}

func (r *WebrtcRepository) GetOfferByRoomID(roomID int) (types.Offer, error) {
	row := r.db.QueryRow("SELECT * FROM offers WHERE room_id = ?", roomID)
	var offer types.Offer
	err := row.Scan(&offer.Offer, &offer.OffererUserID, &offer.RoomID, &offer.OfferIceCandidates, &offer.AnswerIceCandidates)
	if err != nil {
		return types.Offer{}, err
	}	
	return offer, nil
}