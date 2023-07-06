package api

import (
	"movie/graph"
	"movie/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type JSONResponse struct {
	Error   bool        `json:"error"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type MovieGenres struct {
	Movie  *models.Movie   `json:"movie"`
	Genres []*models.Genre `json:"genres"`
}

func (app *Application) Home(c *fiber.Ctx) error {
	var payload = struct {
		Status  string `json:"status"`
		Message string `json:"message"`
		Version string `json:"version"`
	}{
		Status:  "active",
		Message: "Go Movies up and running",
		Version: "1.0.0",
	}
	return c.Status(200).JSON(payload)
}

// AllMovies returns a slice of all movies as JSON.
func (app *Application) AllMovies(c *fiber.Ctx) error {
	movies, err := app.DB.AllMovies()
	if err != nil {
		app.errorJSON(c, err, 400)
		return nil
	}

	return c.Status(200).JSON(movies)
}

func (app *Application) MoviesByGenre(c *fiber.Ctx) error {
	id := c.Params("id")
	genreID, err := strconv.Atoi(id)
	if err != nil {
		app.errorJSON(c, err, 400)
		return nil
	}
	movies, err := app.DB.MoviesByGenreID(genreID)
	if err != nil {
		app.errorJSON(c, err, 400)
		return nil
	}

	return c.Status(200).JSON(movies)
}

func (app *Application) AllGenres(c *fiber.Ctx) error {
	genres, err := app.DB.AllGenres()
	if err != nil {
		app.errorJSON(c, err, 400)
	}
	c.Status(200).JSON(genres)
	return nil
}

func (app *Application) GetAMovie(c *fiber.Ctx) error {
	id := c.Params("id")
	movieID, err := strconv.Atoi(id)
	if err != nil {
		return err
	}

	movie, err := app.getMovieByID(movieID)
	if err != nil {
		return app.errorJSON(c, err, 500)
	}

	return c.Status(200).JSON(movie)
}

func (app *Application) MovieForEdit(c *fiber.Ctx) error {
	id := c.Params("id")
	movieID, err := strconv.Atoi(id)
	if err != nil {
		return err
	}

	movie, err := app.getMovieByID(movieID)
	if err != nil {
		return app.errorJSON(c, err, 500)
	}

	mgIDs, err := app.DB.GetMovieGenresByMovieID(movieID)
	if err != nil {
		return app.errorJSON(c, err, 500)
	}

	genres, err := app.DB.AllGenres()
	if err != nil {
		return app.errorJSON(c, err, 500)
	}

	for ind, genre := range genres {
		if _, ok := mgIDs[genre.ID]; ok {
			genres[ind].Checked = true
		}
	}

	var payload = MovieGenres{
		Movie:  movie,
		Genres: genres,
	}

	return c.Status(200).JSON(payload)
}

func (app *Application) getMovieByID(id int) (*models.Movie, error) {

	movie, err := app.DB.GetAMovie(id)
	if err != nil {
		return nil, err
	}
	return movie, nil
}

func (app *Application) Logout(c *fiber.Ctx) error {
	//clear cache, session if needed
	return c.Next()
}

func (app *Application) AddMove(c *fiber.Ctx) error {
	var movieGenres MovieGenres
	if err := c.BodyParser(&movieGenres); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	cnt, err := app.DB.AddMovie(movieGenres.Movie, movieGenres.Genres)
	if err != nil || cnt == 0 {
		return c.Status(500).JSON(err.Error())
	}

	app.successJSON(c, "movie created")
	return nil
}

func (app *Application) UpdateMovie(c *fiber.Ctx) error {
	var movieGenres MovieGenres
	if err := c.BodyParser(&movieGenres); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	if movieGenres.Movie.ID == 0 {
		id, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return c.Status(400).JSON(err.Error())
		}
		movieGenres.Movie.ID = uint(id)
	}

	cnt, err := app.DB.UpdateMovie(movieGenres.Movie)
	if err != nil || cnt == 0 {
		return c.Status(500).JSON(err.Error())
	}

	err = app.DB.UpdateMovieGenres(movieGenres.Movie.ID, movieGenres.Genres)
	if err != nil {
		return c.Status(500).JSON(err.Error())
	}

	app.successJSON(c, "movie with genres updated")
	return nil
}

// DeleteMovie deletes a movie from the database, by ID.
func (app *Application) DeleteMovie(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(err.Error())
	}

	err = app.DB.DeleteMovie(id)
	if err != nil {
		return c.Status(400).JSON(err.Error())
	}
	app.successJSON(c, "movie deleted")
	return nil
}

func (app *Application) MoviesGraphQL(c *fiber.Ctx) error {
	movies, err := app.DB.AllMovies()
	if err != nil {
		return c.Status(500).JSON(err.Error())
	}

	query := string(c.Body())

	//set the query stirng
	g := graph.New(movies)
	g.QueryString = query

	//perform the query
	resp, err := g.Query()
	if err != nil {
		return c.Status(500).JSON(err.Error())
	}

	return c.Status(200).JSON(resp)
}

// getPoster tries to get a poster image from themoviedb.org.
// func (app *Application) getPoster()  {
// 	type TheMovieDB struct {
// 		Page    int `json:"page"`
// 		Results []struct {
// 			PosterPath string `json:"poster_path"`
// 		} `json:"results"`
// 		TotalPages int `json:"total_pages"`
// 	}

// 	client := &http.Client{}
// 	theUrl := fmt.Sprintf("https://api.themoviedb.org/3/search/movie?api_key=%s", app.APIKey)

// 	// https://api.themoviedb.org/3/search/movie?api_key=b41447e6319d1cd467306735632ba733&query=Die+Hard

// 	req, err := http.NewRequest("GET", theUrl+"&query="+url.QueryEscape(movie.Title), nil)
// 	if err != nil {
// 		log.Println(err)
// 		return movie
// 	}

// 	req.Header.Add("Accept", "application/json")
// 	req.Header.Add("Content-Type", "application/json")

// 	resp, err := client.Do(req)
// 	if err != nil {
// 		log.Println(err)
// 		return movie
// 	}
// 	defer resp.Body.Close()

// 	bodyBytes, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		log.Println(err)
// 		return movie
// 	}

// 	var responseObject TheMovieDB

// 	json.Unmarshal(bodyBytes, &responseObject)

// 	if len(responseObject.Results) > 0 {
// 		movie.Image = responseObject.Results[0].PosterPath
// 	}

// 	return movie
// }

func (app *Application) errorJSON(c *fiber.Ctx, err error, status int) error {
	var payload JSONResponse
	payload.Error = true
	payload.Message = err.Error()

	return c.Status(status).JSON(payload)
}

func (app *Application) successJSON(c *fiber.Ctx, msg string) error {
	var payload JSONResponse
	payload.Error = false
	payload.Message = msg

	return c.Status(200).JSON(payload)
}
