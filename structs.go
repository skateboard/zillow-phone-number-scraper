package main

type response struct {
	AbBucketTreatments any `json:"abBucketTreatments"`
	UserInfo           any `json:"userInfo"`
	PropertyInfo       struct {
		Verified  bool `json:"verified"`
		AgentInfo struct {
			DisplayName          string  `json:"displayName"`
			BusinessName         string  `json:"businessName"`
			PhoneNumber          string  `json:"phoneNumber"`
			AgentBadgeType       string  `json:"agentBadgeType"`
			PhotoURL             string  `json:"photoUrl"`
			ProfileURL           string  `json:"profileUrl"`
			ReviewsReceivedCount int     `json:"reviewsReceivedCount"`
			ReviewsURL           string  `json:"reviewsUrl"`
			RecentSalesCount     any     `json:"recentSalesCount"`
			RatingAverage        float64 `json:"ratingAverage"`
		} `json:"agentInfo"`
		RentalApplicationsEnabled    bool   `json:"rentalApplicationsEnabled"`
		ProviderListingID            any    `json:"providerListingId"`
		ContactFormType              string `json:"contactFormType"`
		MaskType                     string `json:"maskType"`
		MaxLowIncomeList             []any  `json:"maxLowIncomeList"`
		SpecificAvailability         any    `json:"specificAvailability"`
		InstantTourAvailableTimes    any    `json:"instantTourAvailableTimes"`
		InstantTourAvailableTimesMap any    `json:"instantTourAvailableTimesMap"`
		IsLandlordLiaisonProgram     bool   `json:"isLandlordLiaisonProgram"`
		IsIncomeRestricted           bool   `json:"isIncomeRestricted"`
		IncomeRestrictedDisclaimer   any    `json:"incomeRestrictedDisclaimer"`
		IsFeaturedListing            bool   `json:"isFeaturedListing"`
		BuildingName                 any    `json:"buildingName"`
		BuildingAddress              any    `json:"buildingAddress"`
		BestGuessTimeZone            string `json:"bestGuessTimeZone"`
		ZoneID                       any    `json:"zoneId"`
		IsTrustedListing             bool   `json:"isTrustedListing"`
		Reit                         bool   `json:"reit"`
	} `json:"propertyInfo"`
	ContactFormIDsSent any `json:"contactFormIDsSent"`
	RenterProfile      any `json:"renterProfile"`
	InstantTour        any `json:"instantTour"`
}
