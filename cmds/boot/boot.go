package boot

import (
	"encoding/xml"
	"fmt"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
	"gopkg.in/urfave/cli.v1"
	"io/ioutil"
	"net"
	"net/http"
)

type (
	City struct {
		Name string `xml:"Name,attr"`
	}
	State struct {
		Name string `xml:"Name,attr"`
		Code string `xml:"Code,attr"`
		City []City `xml:"City"`
	}
	CountryRegion struct {
		Name  string  `xml:"Name,attr"`
		Code  string  `xml:"Code,attr"`
		State []State `xml:"State"`
	}
	Location struct {
		XMLName       xml.Name `xml:"Location"`
		CountryRegion []CountryRegion
	}
)

func Run(c *cli.Context) (err error) {
	Conf.Load(c)

	//log.Info(v)

	// new echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.Gzip())

	// Routes
	e.GET("/", func(c echo.Context) error {

		type Ret struct {
			Status  int
			Msg     string
			Country []CountryRegion
		}
		r := &Ret{
			Status: 1,
			Msg:    "not found",
		}

		lang := c.QueryParam("lang")
		if lang == "" {
			lang = "zh-cn"
		}

		if lang != "en" && lang != "zh-cn" {
			r.Msg = fmt.Sprintf("lang %s is not supported")
			return c.JSON(http.StatusOK, r)
		}

		raw, err := ioutil.ReadFile(Conf.DataFile + "." + lang)
		if err != nil {
			return err
		}

		location := Location{}
		if err = xml.Unmarshal(raw, &location); err != nil {
			return err
		}

		country := c.QueryParam("country")

		if country == "all" {
			r.Country = location.CountryRegion
			r.Msg = "OK"
			return c.JSON(http.StatusOK, r)
		}

		for _, v := range location.CountryRegion {
			if v.Name == country || v.Code == country {
				r.Status = 0
				r.Msg = "OK"
				r.Country = []CountryRegion{v}
				break
			}
		}
		return c.JSON(http.StatusOK, r)
	})

	// Start server
	l, err := net.Listen("tcp", fmt.Sprintf("%s:%s", Conf.Srv.Host, Conf.Srv.Port))
	if err != nil {
		return err
	}

	e.Listener = l
	e.HideBanner = true

	if Conf.Debug {
		e.Logger.SetLevel(log.DEBUG)
	} else {
		e.Logger.SetLevel(log.ERROR)
	}

	e.Logger.Infof("http server started on %s:%s ", Conf.Srv.Host, Conf.Srv.Port)
	e.Logger.Fatal(e.Start(""))

	return
}
