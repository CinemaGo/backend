package services

import (
	"cinemaGo/backend/internal/models"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type MoviesServiceInterface interface {
	FetchAllCaruselImages() ([]models.CarouselImage, error)
	FetchAllShowsMovie() ([]models.AllShowsMovie, error)
	FetchAShowMovie(showID int) (models.AShowMovie, error)
	FetchAllActorsCrewsByMovieID(movieID int) ([]models.ActorsCrewsOfMovie, error)
	FetchActorCrewInfo(actorCrewID int) (models.ActorCrewInfo, error)
	FetchMoviesByActorCrewID(actorCrewID int) ([]models.ActorCrewMovies, error)
}

type MoviesService struct {
	db          models.DBContractMovies
	redisClient *redis.Client
}

func NewMoviesService(db models.DBContractMovies) (*MoviesService, error) {

	redisAddr, redisPass, err := LoadRedisEnvironmentVariables("REDIS_ADDR", "REDIS_PASS")
	if err != nil {
		return nil, err
	}
	rdb := redis.NewClient(&redis.Options{
		Addr:        redisAddr,
		Password:    redisPass,
		DB:          1,
		DialTimeout: 5 * time.Second,
	})

	return &MoviesService{
		db:          db,
		redisClient: rdb,
	}, nil
}

// FetchAllCaruselImages fetches carousel images, first checking the Redis cache for existing data.
// If the data is not found in the cache, it retrieves the data from the database and stores it in the cache.
//
// Parameters:
// - None
//
// Returns:
// - []models.CarouselImage: A slice containing the carousel images.
// - error: If an error occurs during the fetching process from either Redis or the database.
func (ms *MoviesService) FetchAllCaruselImages() ([]models.CarouselImage, error) {
	carouselImages, err := ms.fetchCachedCarouselImagesDataFromRedis("carouselImages")
	if errors.Is(err, ErrNoCachedDataFound) {
		carouselImages, err := ms.db.RetrieveAllCarouselImages()
		if err != nil {
			return nil, fmt.Errorf("error occurred while fetching all carousel images in the service section: %w", err)
		}

		if err := ms.cacheCarouselImagesDataInRedisCache(carouselImages); err != nil {
			return nil, err
		}

		return carouselImages, nil
	}

	if err != nil {
		return nil, err
	}

	return carouselImages, nil
}

// cacheCarouselImagesDataInRedisCache stores carousel images data in Redis cache for future retrieval.
//
// Parameters:
// - carouselImages []models.CarouselImage: The slice of carousel images that will be cached in Redis.
//
// Returns:
// - error: Returns an error if any error occurs during the marshalling of the data or setting it in Redis.
func (ms *MoviesService) cacheCarouselImagesDataInRedisCache(carouselImages []models.CarouselImage) error {
	jsonData, err := json.Marshal(carouselImages)
	if err != nil {
		return fmt.Errorf("error occurred while marshalling carousel images, to cache for Redis: %w", err)
	}

	err = ms.redisClient.Set(context.Background(), "carouselImages", jsonData, 30*time.Minute).Err()
	if err != nil {
		return fmt.Errorf("error occurred while setting up carouselImages data in Redis: %w", err)
	}

	return nil
}

// fetchCachedCarouselImagesDataFromRedis fetches carousel images data from Redis cache using the provided key.
//
// Parameters:
// - keyName string: The key in Redis to retrieve the carousel images data.
//
// Returns:
// - []models.CarouselImage: A slice of carousel images data retrieved from the Redis cache.
// - error: Returns an error if the data is not found or an issue occurs while retrieving or unmarshalling the data.
func (ms *MoviesService) fetchCachedCarouselImagesDataFromRedis(keyName string) ([]models.CarouselImage, error) {
	data, err := ms.redisClient.Get(context.Background(), keyName).Result()
	if errors.Is(err, redis.Nil) {
		return nil, ErrNoCachedDataFound
	}

	if err != nil {
		return nil, fmt.Errorf("error occurred while fetching cached carouselImages data from Redis: %w", err)
	}

	var carouselImagesData []models.CarouselImage

	err = json.Unmarshal([]byte(data), &carouselImagesData)
	if err != nil {
		return nil, fmt.Errorf("error occurred while unmarshalling carouselImages data that is coming from Redis cache: %w", err)
	}

	return carouselImagesData, nil
}

// FetchAllShowsMovie retrieves all show movies, first checking the Redis cache for existing data.
// If the data is not found in the cache, it fetches the data from the database, processes it,
// and stores it in Redis cache for future requests.
//
// Parameters:
// - None
//
// Returns:
// - []models.AllShowsMovie: A slice containing all show movies with their respective ratings.
// - error: If an error occurs during the fetching process from Redis or the database.
func (ms *MoviesService) FetchAllShowsMovie() ([]models.AllShowsMovie, error) {

	showsMovie, err := ms.fetchAllShowsMovieDataFromRedisCache("showsMovie")
	if errors.Is(err, ErrNoCachedDataFound) {
		showsMovie, err := ms.db.RetrieveAllShowsMovie()
		if err != nil {
			return nil, fmt.Errorf("error occurred while fetching all showsMovie in the service section: %w", err)
		}

		for i := range showsMovie {
			showsMovie[i].MovieRating = showsMovie[i].MovieRating / 10
		}

		if err = ms.cacheAllShowsMovieDataInRedis(showsMovie); err != nil {
			return nil, err
		}

		return showsMovie, nil
	}

	if err != nil {
		return nil, err
	}

	return showsMovie, nil
}

// cacheAllShowsMovieDataInRedis caches all show movie data in Redis for future use.
//
// Parameters:
// - showsMovie []models.AllShowsMovie: The slice of show movies to be cached in Redis.
//
// Returns:
// - error: Returns an error if there is an issue during the marshalling or setting of the data in Redis.
func (ms *MoviesService) cacheAllShowsMovieDataInRedis(showsMovie []models.AllShowsMovie) error {
	jsonData, err := json.Marshal(showsMovie)
	if err != nil {
		return fmt.Errorf("error occurred while marshalling showsMovie, to cache for Redis: %w", err)
	}

	err = ms.redisClient.Set(context.Background(), "showsMovie", jsonData, 30*time.Minute).Err()
	if err != nil {
		return fmt.Errorf("error occurred while setting up showsMovie data in Redis: %w", err)
	}

	return nil
}

// fetchAllShowsMovieDataFromRedisCache fetches the show movies data from Redis cache using the provided key.
//
// Parameters:
// - keyName string: The Redis key under which the show movies data is cached.
//
// Returns:
// - []models.AllShowsMovie: A slice of show movie data retrieved from Redis cache.
// - error: Returns an error if the data is not found in the cache, or if there is an issue retrieving or unmarshalling the data.
func (ms *MoviesService) fetchAllShowsMovieDataFromRedisCache(keyName string) ([]models.AllShowsMovie, error) {
	data, err := ms.redisClient.Get(context.Background(), keyName).Result()
	if errors.Is(err, redis.Nil) {
		return nil, ErrNoCachedDataFound
	}

	if err != nil {
		return nil, fmt.Errorf("error occurred while fetching cached showsMovie data from Redis: %w", err)
	}

	var showsMovieData []models.AllShowsMovie

	err = json.Unmarshal([]byte(data), &showsMovieData)
	if err != nil {
		return nil, fmt.Errorf("error occurred while unmarshalling showsMovie data that is coming from Redis cache: %w", err)
	}

	return showsMovieData, nil
}

// FetchAShowMovie retrieves a specific movie by its ID. First, it checks if the data is available in Redis cache.
// If not found in cache, it fetches the data from the database, processes it, and caches the data in Redis for future requests.
//
// Parameters:
// - showID int: The ID of the movie to retrieve.
//
// Returns:
// - models.AShowMovie: A Movie object containing the movie details.
// - error: Returns an error if the movie is not found or if there is an issue fetching the movie from the database or cache.
func (ms *MoviesService) FetchAShowMovie(showID int) (models.AShowMovie, error) {
	aShowMovie, err := ms.fetchAShowMovieDataFromRedisCache(fmt.Sprintf("%v", showID))
	if errors.Is(err, ErrNoCachedDataFound) {
		aShowMovie, err := ms.db.RetrieveAShowMovie(showID)
		if err != nil {
			if errors.Is(err, models.ErrMovieNotFoundByID) {
				return models.AShowMovie{}, ErrMovieNotFoundByID
			}

			return models.AShowMovie{}, fmt.Errorf("error occurred while fetching a aShowMovie by id in the service section: %w", err)
		}

		aShowMovie.MovieRating = float32(aShowMovie.MovieRating) / 10.0

		if err := ms.cacheAShowMovieDataInRedis(aShowMovie); err != nil {
			return models.AShowMovie{}, err
		}

		return aShowMovie, nil
	}

	if err != nil {
		return models.AShowMovie{}, nil
	}

	return aShowMovie, nil
}

// cacheAShowMovieDataInRedis caches a single movie data in Redis for future use.
//
// Parameters:
// - aShowMovie models.AShowMovie: The movie data to be cached in Redis.
//
// Returns:
// - error: Returns an error if there is an issue during the marshalling or setting of the data in Redis.
func (ms *MoviesService) cacheAShowMovieDataInRedis(aShowMovie models.AShowMovie) error {
	jsonData, err := json.Marshal(aShowMovie)
	if err != nil {
		return fmt.Errorf("error occurred while marshalling aShowMovie, to cache for Redis: %w", err)
	}

	err = ms.redisClient.Set(context.Background(), fmt.Sprintf("%v", aShowMovie.ShowID), jsonData, 30*time.Minute).Err()
	if err != nil {
		return fmt.Errorf("error occurred while setting up aShowMovie data in Redis: %w", err)
	}

	return nil
}

// fetchAShowMovieDataFromRedisCache fetches the movie data from Redis cache using the provided key.
//
// Parameters:
// - keyName string: The Redis key under which the movie data is cached.
//
// Returns:
// - models.AShowMovie: The movie data retrieved from Redis cache.
// - error: Returns an error if the data is not found in the cache, or if there is an issue retrieving or unmarshalling the data.
func (ms *MoviesService) fetchAShowMovieDataFromRedisCache(keyName string) (models.AShowMovie, error) {
	data, err := ms.redisClient.Get(context.Background(), keyName).Result()
	if errors.Is(err, redis.Nil) {
		return models.AShowMovie{}, ErrNoCachedDataFound
	}

	if err != nil {
		return models.AShowMovie{}, fmt.Errorf("error occurred while fetching cached aShowMovie data from Redis: %w", err)
	}

	var aShowMovie models.AShowMovie

	err = json.Unmarshal([]byte(data), &aShowMovie)
	if err != nil {
		return models.AShowMovie{}, fmt.Errorf("error occurred while unmarshalling aShowMovie data that is coming from Redis cache: %w", err)
	}

	return aShowMovie, nil
}

// FetchAllActorsCrewsByMovieID retrieves a list of actors and crew for a specific movie by its ID.
// It first attempts to fetch the data from Redis cache, and if not found, retrieves the data from the database
// and caches it in Redis for future use.
//
// Parameters:
// - movieID int: The ID of the movie for which to fetch the actors and crew.
//
// Returns:
// - []models.ActorsCrewsOfMovie: A list of actors and crew associated with the specified movie.
// - error: Returns an error if the data cannot be fetched from the cache or database, or if there is an issue with caching.
func (ms *MoviesService) FetchAllActorsCrewsByMovieID(movieID int) ([]models.ActorsCrewsOfMovie, error) {
	allActorsCrew, err := ms.fetchAllActorsCrewsByMovieIDFromRedisCache(fmt.Sprintf("%v", movieID))
	if errors.Is(err, ErrNoCachedDataFound) {
		allActorsCrew, err := ms.db.RetrieveAllActorsCrewsByMovieID(movieID)
		if err != nil {
			return nil, fmt.Errorf("error occurred while fetching all actors and crews in the service section: %w", err)
		}

		if err := ms.cacheAllActorsCrewsByMovieIDInRedis(movieID, allActorsCrew); err != nil {
			return nil, err
		}

		return allActorsCrew, nil
	}

	if err != nil {
		return nil, err
	}

	return allActorsCrew, nil
}

// cacheAllActorsCrewsByMovieIDInRedis caches a list of actors and crew data for a specific movie in Redis.
//
// Parameters:
// - allActorsCrew []models.ActorsCrewsOfMovie: The list of actors and crew data to be cached.
//
// Returns:
// - error: Returns an error if there is an issue during the marshalling or setting of the data in Redis.
func (ms *MoviesService) cacheAllActorsCrewsByMovieIDInRedis(movieID int, allActorsCrew []models.ActorsCrewsOfMovie) error {

	jsonData, err := json.Marshal(allActorsCrew)
	if err != nil {
		return fmt.Errorf("error occurred while marshalling allActorsCrew, to cache for Redis: %w", err)
	}

	err = ms.redisClient.Set(context.Background(), fmt.Sprintf("allActorsCrew%v", movieID), jsonData, 30*time.Minute).Err()
	if err != nil {
		return fmt.Errorf("error occurred while setting up allActorsCrew data in Redis: %w", err)
	}

	return nil
}

// fetchAllActorsCrewsByMovieIDFromRedisCache retrieves the list of actors and crew data from Redis cache
// for a specific movie using the provided cache key.
//
// Parameters:
// - keyName string: The Redis key under which the actors and crew data is cached.
//
// Returns:
// - []models.ActorsCrewsOfMovie: A list of actors and crew for the movie.
// - error: Returns an error if the data is not found in the cache, or if there is an issue retrieving or unmarshalling the data.
func (ms *MoviesService) fetchAllActorsCrewsByMovieIDFromRedisCache(keyName string) ([]models.ActorsCrewsOfMovie, error) {
	data, err := ms.redisClient.Get(context.Background(), fmt.Sprintf("allActorsCrew%v", keyName)).Result()
	if errors.Is(err, redis.Nil) {
		return nil, ErrNoCachedDataFound
	}

	if err != nil {
		return nil, fmt.Errorf("error occurred while fetching cached allActorsCrew data from Redis: %w", err)
	}

	var allActorsCrew []models.ActorsCrewsOfMovie

	err = json.Unmarshal([]byte(data), &allActorsCrew)
	if err != nil {
		return nil, fmt.Errorf("error occurred while unmarshalling allActorsCrew data that is coming from Redis cache: %w", err)
	}

	return allActorsCrew, nil
}

// FetchActorCrewInfo retrieves detailed information about a specific actor or crew member by their ID.
// It first tries to fetch the data from Redis cache, and if not found, fetches the data from the database
// and caches it for future use.
//
// Parameters:
// - actorCrewID int: The ID of the actor or crew member whose information is being fetched.
//
// Returns:
// - models.ActorCrewInfo: The detailed information about the actor or crew member.
// - error: Returns an error if the data cannot be fetched from the cache or database, or if there is an issue with caching.
func (ms *MoviesService) FetchActorCrewInfo(actorCrewID int) (models.ActorCrewInfo, error) {
	actorCrewInfo, err := ms.fetchActorCrewInfoFromRedisCache(fmt.Sprintf("%v", actorCrewID))
	if errors.Is(err, ErrNoCachedDataFound) {
		actorCrewInfo, err := ms.db.RetriveActorCrewInfo(actorCrewID)
		if err != nil {
			if errors.Is(err, models.ErrActorCrewNotFoundByID) {
				return models.ActorCrewInfo{}, ErrActorCrewNotFoundByID
			}
			return models.ActorCrewInfo{}, fmt.Errorf("error occurred while fetching actor or crew info in the service section: %w", err)
		}

		if err := ms.cacheActorCrewInfoIDInRedis(actorCrewInfo); err != nil {
			return models.ActorCrewInfo{}, err
		}

		return actorCrewInfo, nil
	}

	if err != nil {
		return models.ActorCrewInfo{}, err
	}

	return actorCrewInfo, nil
}

// cacheActorCrewInfoIDInRedis caches the actor or crew information in Redis for a specified period of time.
//
// Parameters:
// - ActorCrewInfo models.ActorCrewInfo: The actor or crew member's information to cache.
//
// Returns:
// - error: Returns an error if there is an issue during the marshalling or setting of the data in Redis.
func (ms *MoviesService) cacheActorCrewInfoIDInRedis(ActorCrewInfo models.ActorCrewInfo) error {
	jsonData, err := json.Marshal(ActorCrewInfo)
	if err != nil {
		return fmt.Errorf("error occurred while marshalling ActorCrewInfo, to cache for Redis: %w", err)
	}

	err = ms.redisClient.Set(context.Background(), fmt.Sprintf("actorCrewID%v", ActorCrewInfo.ID), jsonData, 30*time.Minute).Err()
	if err != nil {
		return fmt.Errorf("error occurred while setting up ActorCrewInfo data in Redis: %w", err)
	}

	return nil
}

// fetchActorCrewInfoFromRedisCache retrieves the actor or crew information from Redis cache using the provided key.
//
// Parameters:
// - keyName string: The Redis key under which the actor or crew member's information is cached.
//
// Returns:
// - models.ActorCrewInfo: The detailed information about the actor or crew member.
// - error: Returns an error if the data is not found in the cache, or if there is an issue retrieving or unmarshalling the data.
func (ms *MoviesService) fetchActorCrewInfoFromRedisCache(keyName string) (models.ActorCrewInfo, error) {
	data, err := ms.redisClient.Get(context.Background(), fmt.Sprintf("actorCrewID%v", keyName)).Result()
	if errors.Is(err, redis.Nil) {
		return models.ActorCrewInfo{}, ErrNoCachedDataFound
	}

	if err != nil {
		return models.ActorCrewInfo{}, fmt.Errorf("error occurred while fetching cached ActorCrewInfo data from Redis: %w", err)
	}
	var ActorCrewInfo models.ActorCrewInfo

	err = json.Unmarshal([]byte(data), &ActorCrewInfo)
	if err != nil {
		return models.ActorCrewInfo{}, fmt.Errorf("error occurred while unmarshalling ActorCrewInfo data that is coming from Redis cache: %w", err)
	}

	return ActorCrewInfo, nil
}

// FetchMoviesByActorCrewID retrieves a list of movies associated with a specific actor or crew member
// using their actor/crew ID. It first attempts to fetch the data from the Redis cache, and if not found,
// retrieves it from the database and caches the result for future use.
//
// Parameters:
// - actorCrewID int: The ID of the actor or crew member whose movies are being fetched.
//
// Returns:
// - []models.ActorCrewMovies: A list of movies associated with the actor or crew member.
// - error: Returns an error if the data cannot be fetched from the cache or database, or if there is an issue with caching.
func (ms *MoviesService) FetchMoviesByActorCrewID(actorCrewID int) ([]models.ActorCrewMovies, error) {
	actorCrewMovies, err := ms.fetchMoviesByActorCrewIDFromRedisCache(fmt.Sprintf("%v", actorCrewID))
	if errors.Is(err, ErrNoCachedDataFound) {
		actorCrewMovies, err := ms.db.RetrieveMoviesByActorCrewID(actorCrewID)
		if err != nil {
			if errors.Is(err, models.ErrActorCrewNotFoundByID) {
				return nil, ErrActorCrewNotFoundByID
			}
			return nil, fmt.Errorf("error occurred while fetching all movies of actors or crews in the service section: %w", err)
		}

		if err := ms.cacheMoviesByActorCrewIDInRedis(actorCrewID, actorCrewMovies); err != nil {
			return nil, err
		}

		return actorCrewMovies, nil
	}

	if err != nil {
		return nil, err
	}

	return actorCrewMovies, nil
}

// cacheMoviesByActorCrewIDInRedis caches the list of movies associated with a specific actor or crew
// in Redis for a specified period of time.
//
// Parameters:
// - actorCrewMovies []models.ActorCrewMovies: A list of movies associated with the actor or crew to cache.
//
// Returns:
// - error: Returns an error if there is an issue during the marshalling or setting of the data in Redis.
func (ms *MoviesService) cacheMoviesByActorCrewIDInRedis(actorCrewID int, ActorCrewMovies []models.ActorCrewMovies) error {
	jsonData, err := json.Marshal(ActorCrewMovies)
	if err != nil {
		return fmt.Errorf("error occurred while marshalling actorCrewMovies, to cache for Redis: %w", err)
	}

	err = ms.redisClient.Set(context.Background(), fmt.Sprintf("MoviesByActorCrewID%v", actorCrewID), jsonData, 30*time.Minute).Err()
	if err != nil {
		return fmt.Errorf("error occurred while setting up actorCrewMovies data in Redis: %w", err)
	}

	return nil
}

// fetchMoviesByActorCrewIDFromRedisCache retrieves the list of movies associated with a specific actor
// or crew member from Redis cache using the provided key.
//
// Parameters:
// - keyName string: The Redis key under which the actor or crew member's movies are cached.
//
// Returns:
// - []models.ActorCrewMovies: A list of movies associated with the actor or crew member.
// - error: Returns an error if the data is not found in the cache, or if there is an issue retrieving or unmarshalling the data.
func (ms *MoviesService) fetchMoviesByActorCrewIDFromRedisCache(keyName string) ([]models.ActorCrewMovies, error) {
	data, err := ms.redisClient.Get(context.Background(), fmt.Sprintf("MoviesByActorCrewID%v", keyName)).Result()
	if errors.Is(err, redis.Nil) {
		return nil, ErrNoCachedDataFound
	}

	if err != nil {
		return nil, fmt.Errorf("error occurred while fetching cached actorCrewMovies data from Redis: %w", err)
	}
	var actorCrewMovies []models.ActorCrewMovies

	err = json.Unmarshal([]byte(data), &actorCrewMovies)
	if err != nil {
		return nil, fmt.Errorf("error occurred while unmarshalling actorCrewMovies data that is coming from Redis cache: %w", err)
	}

	return actorCrewMovies, nil
}
