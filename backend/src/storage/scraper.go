package storage

import (
	"database/sql"
	"fmt"
	"search_scraper/src/types"
	"search_scraper/src/utils"
)

type ScrapedResult struct {
	ScrapedLinks []types.ScrapedLink `json:"scraped_link"`
	ResultRating float64             `json:"result_rating"`
}

func (s *Storage) FilterLinks(gsl []types.ScrapedLink, p types.ScrapedPage) []types.ScrapedLink {
	fmt.Println(p.NextPageLink)
	var sl []types.ScrapedLink
	for _, l := range p.ScrapedLink {
		wl, err := s.ConteinsLinkInList("whitelist", types.Link{
			Url:    l.Link,
			Domain: l.Domain,
		})
		if err != nil && err != sql.ErrNoRows {
			fmt.Println(err.Error())
		}
		bl, err := s.ConteinsLinkInList("blacklist", types.Link{
			Url:    l.Link,
			Domain: l.Domain,
		})
		if err != nil && err != sql.ErrNoRows {
			fmt.Println(err)
		}
		fl, err := s.ConteinsLinkInList("findedlist", types.Link{
			Url:    l.Link,
			Domain: l.Domain,
		})
		if err != nil && err != sql.ErrNoRows {
			fmt.Println(err)
		}

		for _, gl := range gsl {
			if l.Domain == gl.Domain {
				fl = true
				break
			}
		}

		if (!fl && !bl) || (wl && !bl) {
			sl = append(sl, l)
			fmt.Printf("\t\t%30s\t%t\t%t\t%t\n", l.Domain, fl, wl, bl)
		}
	}
	return sl

}

func (s *Storage) FilteredScraping(q string, d int) (ScrapedResult, error) {
	var sl []types.ScrapedLink
	fp, err := utils.Srcape(fmt.Sprintf("/search?q=%s", q))
	if err != nil {
		fmt.Println(err)
		return ScrapedResult{sl, 0}, err
	}
	sl = append(sl, s.FilterLinks(sl, fp)...)

	total := len(fp.ScrapedLink)

	for i := 0; i < d && fp.NextPageLink != ""; i++ {
		fmt.Println(i)
		p, err := utils.Srcape(fp.NextPageLink)
		if err != nil {
			fmt.Println(err)
			return ScrapedResult{sl, 0}, err
		}
		total += len(p.ScrapedLink)
		sl = append(sl, s.FilterLinks(sl, p)...)
	}
	fmt.Printf("T: %d, F: %d (%2.2f)", total, len(sl), float64(len(sl))/float64(total)*100)
	r := float64(len(sl)) / float64(total) * 100
	return ScrapedResult{sl, r}, nil
}