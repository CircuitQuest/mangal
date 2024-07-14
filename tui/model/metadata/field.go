package metadata

import "github.com/luevano/libmangal/metadata"

const (
	titleSize       = 80
	scoreSize       = 3
	descriptionSize = 80
	tagsSize        = 20
	peopleSize      = 20
	urlSize         = 100
	dateSize        = 10
	publisherSize   = 10 // need to tweak
	statusSize      = 10
	formatSize      = 10
	countrySize     = 2
	chaptersSize    = 5
	idSize          = 6
	extraIDsSize    = 40
)

// field represents a key-value pair with
// a given max width for the rendered value
//
// the minimum width is 4; for a list, the width
// applies only on the list items
type field struct {
	name  string
	value any
	width int
}

// fields represents all possible metadata fields
type fields struct {
	title,
	alternateTitles,
	score,
	description,
	cover,
	banner,
	tags,
	genres,
	characters,
	authors,
	artists,
	translators,
	letterers,
	startDate,
	endDate,
	publisher,
	status,
	format,
	country,
	chapters,
	notes,
	url,
	id,
	extraIDs field
}

func allFields(meta metadata.Metadata) fields {
	return fields{
		title:           field{"Title", meta.Title(), titleSize},
		alternateTitles: field{"Alternate titles", meta.AlternateTitles(), titleSize},
		score:           field{"Score", meta.Score(), scoreSize},
		description:     field{"Description", meta.Description(), descriptionSize},
		cover:           field{"Cover", meta.Cover(), urlSize},
		banner:          field{"Banner", meta.Banner(), urlSize},
		tags:            field{"Tags", meta.Tags(), tagsSize},
		genres:          field{"Genres", meta.Genres(), tagsSize},
		characters:      field{"Characters", meta.Characters(), peopleSize},
		authors:         field{"Authors", meta.Authors(), peopleSize},
		artists:         field{"Artists", meta.Artists(), peopleSize},
		translators:     field{"Translators", meta.Translators(), peopleSize},
		letterers:       field{"Letterers", meta.Letterers(), peopleSize},
		startDate:       field{"Start date", meta.StartDate(), dateSize},
		endDate:         field{"End date", meta.EndDate(), dateSize},
		publisher:       field{"Publisher", meta.Publisher(), publisherSize},
		status:          field{"Status", meta.Status(), statusSize},
		format:          field{"Format", meta.Format(), formatSize},
		country:         field{"Country", meta.Country(), countrySize},
		chapters:        field{"Chapters", meta.Chapters(), chaptersSize},
		notes:           field{"Notes", meta.Notes(), descriptionSize},
		url:             field{"URL", meta.URL(), urlSize},
		id:              field{"ID", meta.ID(), idSize},
		extraIDs:        field{"External IDs", meta.ExtraIDs(), extraIDsSize},
	}
}
