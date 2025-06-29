package types

type WebrtcRepository interface {
	CreateOffer(offer Offer) (Offer, error)
	GetOfferByRoomID(roomID int) (Offer, error)
}

