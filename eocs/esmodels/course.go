package esmodels

import (
	"time"
)

type Course struct {
	ID                 string            `bson:"_id"`
	IsOrganizationOnly bool              `bson:"is_organization_only"`
	Headline           IntlStringWrapper `bson:"headline"`
	Description        IntlStringWrapper `bson:"description"`
	SubscriptionLevel  int               `bson:"subscription_level"`
	ViewCount          int               `bson:"view_count"`
	Title              IntlStringWrapper `bson:"title"`
	EnrolledCount      int               `bson:"enrolled_count"`
	SkillLevel         int               `bson:"skill_level"`
	EstMinutes         int               `bson:"est_minutes"`
	PrimaryTopic       string            `bson:"primary_topic"`
	Units              UnitsWrapper      `bson:"units"`
	CoverURL           string            `bson:"cover_url"`
	LogoURL            string            `bson:"logo_url"`
	IsPublished        bool              `bson:"is_published"`
	InfoMD             string            `bson:"info_md"`
	VerifiedCertCost   float64           `bson:"verified_cert_cost"`
	OrganizationIDs    []string          `bson:"organization_ids"`
	Topics             []string          `bson:"topics"`
	CreatedAt          time.Time         `bson:"created_at"`
	UpdatedAt          time.Time         `bson:"updated_at"`
}
