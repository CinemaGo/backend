package handlers

import (
	"cinemaGo/backend/api/helpers"
	"cinemaGo/backend/internal/services"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type BookingHandler struct {
	booking services.BookingServiceInterface
}

func NewBookingHandler(service services.BookingServiceInterface) *BookingHandler {
	return &BookingHandler{booking: service}
}

func (service *BookingHandler) MovieShowTimes(c *gin.Context) {
	showID, err := helpers.GetParameterFromURL(c, "showID", "invalid show ID provided.")
	if err != nil {
		helpers.ClientError(c, http.StatusBadRequest, fmt.Sprintf("%v", err))
		return
	}

	showMovieInfo, err := service.booking.FetchShowMovieInfo(showID)
	if err != nil {
		if errors.Is(err, services.ErrShowNotFound) {
			helpers.ClientError(c, http.StatusNotFound, fmt.Sprintf("show not found by %v ID", showID))
			return
		}
		helpers.ServerError(c, err)
		return
	}

	showInfos, err := service.booking.FetchShowInfo(showMovieInfo.MovieID)
	if err != nil {
		if errors.Is(err, services.ErrShowNotFound) {
			helpers.ClientError(c, http.StatusNotFound, fmt.Sprintf("show not found by %v ID", showID))
			return
		}
		helpers.ServerError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"showMovieInfo": showMovieInfo,
		"showInfos":     showInfos,
	})
}

func (service *BookingHandler) ShowSeats(c *gin.Context) {
	showID, err := helpers.GetParameterFromURL(c, "showID", "invalid show ID provided.")
	if err != nil {
		helpers.ClientError(c, http.StatusBadRequest, fmt.Sprintf("%v", err))
		return
	}

	showSeats, err := service.booking.FetchShowSeats(showID)
	if err != nil {
		if errors.Is(err, services.ErrShowSeatNotFound) {
			helpers.ClientError(c, http.StatusNotFound, fmt.Sprintf("show ID %v not found", showID))
			return
		}
		helpers.ServerError(c, err)
		return
	}

	showInfo, err := service.booking.FetchShowSeatsMovieInfo(showID)
	if err != nil {
		if errors.Is(err, services.ErrShowNotFound) {
			helpers.ClientError(c, http.StatusNotFound, fmt.Sprintf("show ID %v not found", showID))
			return
		}
		helpers.ServerError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"showSeats": showSeats,
		"showInfo":  showInfo,
	})
}

func (service *BookingHandler) BookSeats(c *gin.Context) {
	var bookingForm BookingForm

	if err := c.ShouldBindJSON(&bookingForm); err != nil {
		helpers.RespondWithValidationErrors(c, err, bookingForm)
		return
	}

	if err := helpers.ValidateSeatsID(bookingForm.ShowSeatsID); err != nil {
		helpers.ClientError(c, http.StatusBadRequest, fmt.Sprintf("%v", err))
		return
	}

	userID, _ := c.Get("userID")
	user_id := int(userID.(float64))

	err := service.booking.CreateNewBooking(bookingForm.ShowID, user_id, bookingForm.ShowSeatsID)
	if err != nil {
		if errors.Is(err, services.ErrShowNotFound) {
			helpers.ClientError(c, http.StatusNotFound, fmt.Sprintf("show ID %v not found", bookingForm.ShowID))
			return
		}

		if errors.Is(err, services.ErrShowSeatHasSelected) {
			helpers.ClientError(c, http.StatusConflict, "Sorry! These seats are no longer available. Please try again with other seats.")
			return
		}

		if errors.Is(err, services.ErrTooManySeats) {
			helpers.ClientError(c, http.StatusBadRequest, "You can select a maximum of 5 seats at a time.")
			return
		}

		helpers.ServerError(c, err)
		return
	}

	c.JSON(http.StatusOK,gin.H{
		"message":"Booking successful! Your seats are reserved, and payment has been completed. Enjoy the show!",
	})
}
