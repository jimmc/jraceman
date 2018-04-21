package domain

type Repos interface {
  Area() AreaRepo
  Site() SiteRepo
}
