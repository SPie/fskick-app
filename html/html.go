package html

import (
	"io"
	"text/template"

	"github.com/spie/fskick/games"
)

type HtmlWriter interface {
	WriteSeasonTable(season games.Season, head []string, rows [][]string) error
}

type htmlWriter struct {
	writer io.Writer
}

func NewHtmlWriter(writer io.Writer) HtmlWriter {
	return htmlWriter{writer: writer}
}

func (htmlWriter htmlWriter) WriteSeasonTable(season games.Season, head []string, rows [][]string) error {
	template, err := template.New("SeasonTable").Parse(getSeasonTableTemplate())
	if err != nil {
		return err
	}

	err = template.ExecuteTemplate(htmlWriter.writer, "SeasonTable", seasonTableData{
		SeasonName: season.Name,
		TableHead:  head,
		TableRows:  rows,
	})
	if err != nil {
		return err
	}

	return nil
}

type seasonTableData struct {
	SeasonName string
	TableHead  []string
	TableRows  [][]string
}

func getSeasonTableTemplate() string {
	return `<?DOCTYPE html>
<html>
    <head>
        <title>FS Kick Table {{.SeasonName}}</title>
        <link rel="stylesheet" type="text/css" href="styles.css" />
    </head>
    <body>
        <h1>{{.SeasonName}}</h1>
        <table>
            <tr>
                {{range $head := .TableHead}}<th>{{$head}}</th>
				{{end}} 
            </tr>
            {{range $row := .TableRows}}
			<tr>
				{{range $col := $row}}<td>{{$col}}</td>
				{{end}}
			</tr>{{end}}
        </table>
    </body>
</html>`
}
