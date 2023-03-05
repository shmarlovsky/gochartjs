package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	g "github.com/maragudk/gomponents"
	c "github.com/maragudk/gomponents/components"
	. "github.com/maragudk/gomponents/html"
)

type AxisData []string

// print in JS style
func (d AxisData) String() string {
	s := "["
	for _, e := range d {
		s += fmt.Sprintf(" '%s',", e)
	}
	s += "]"
	return s
}

type PlotData struct {
	X AxisData
	Y AxisData
}

var plotData PlotData = PlotData{
	X: AxisData{"Red", "Blue", "Yellow", "Green", "Purple", "Orange"},
	Y: AxisData{"12", "19", "3", "5", "2", "3"},
}

func main() {

	http.Handle("/", createHandler(Chart()))

	if err := http.ListenAndServe("localhost:8081", nil); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Println("Error:", err)
	}
}

func createHandler(title string, body g.Node) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Rendering a Node is as simple as calling Render and passing an io.Writer
		_ = Page(title, r.URL.Path, body).Render(w)
	}
}

func Page(title, path string, body g.Node) g.Node {
	// HTML5 boilerplate document
	return c.HTML5(c.HTML5Props{
		Title:    title,
		Language: "en",
		Body:     []g.Node{body},
	})
}

func Chart() (string, g.Node) {
	return "charts.js wrapper",
		Div(
			chartJsCdnScript(),
			chartCanvasDiv("plot1"),
			barChartScript("plot1", plotData.X, plotData.Y),
		)

}

func chartJsCdnScript() g.Node {
	return Script(Src("https://cdn.jsdelivr.net/npm/chart.js"))
}

func chartCanvasDiv(id string) g.Node {
	return Div(
		g.Raw(fmt.Sprintf(`
	<div>
		<canvas id="%s"></canvas>
  	</div>`, id),
		),
	)
}

func barChartScript(id string, x AxisData, y AxisData) g.Node {
	return g.Raw(
		fmt.Sprintf(`<script>
		const ctx = document.getElementById('%s');
	  
		new Chart(ctx, {
		  type: 'bar',
		  data: {
			labels: %v,
			datasets: [{
			  label: '# of Votes',
			  data: %v,
			  borderWidth: 1
			}]
		  },
		  options: {
			scales: {
			  y: {
				beginAtZero: true
			  }
			}
		  }
		});
	  </script>`, id, x, y),
	)
}
