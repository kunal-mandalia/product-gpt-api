package search

type EbaySearchResponse struct {
	AutoCorrections struct {
		Q string `json:"q"`
	} `json:"autoCorrections"`
	Href          string `json:"href"`
	ItemSummaries []struct {
		AdditionalImages []struct {
			Height   string `json:"height"`
			ImageURL string `json:"imageUrl"`
			Width    string `json:"width"`
		} `json:"additionalImages"`
		AdultOnly        string   `json:"adultOnly"`
		AvailableCoupons string   `json:"availableCoupons"`
		BidCount         string   `json:"bidCount"`
		BuyingOptions    []string `json:"buyingOptions"`
		Categories       []struct {
			CategoryID   string `json:"categoryId"`
			CategoryName string `json:"categoryName"`
		} `json:"categories"`
		CompatibilityMatch      string `json:"compatibilityMatch"`
		CompatibilityProperties []struct {
			LocalizedName string `json:"localizedName"`
			Name          string `json:"name"`
			Value         string `json:"value"`
		} `json:"compatibilityProperties"`
		Condition       string `json:"condition"`
		ConditionID     string `json:"conditionId"`
		CurrentBidPrice struct {
			ConvertedFromCurrency string `json:"convertedFromCurrency"`
			ConvertedFromValue    string `json:"convertedFromValue"`
			Currency              string `json:"currency"`
			Value                 string `json:"value"`
		} `json:"currentBidPrice"`
		DistanceFromPickupLocation struct {
			UnitOfMeasure string `json:"unitOfMeasure"`
			Value         string `json:"value"`
		} `json:"distanceFromPickupLocation"`
		EnergyEfficiencyClass string `json:"energyEfficiencyClass"`
		Epid                  string `json:"epid"`
		Image                 struct {
			Height   string `json:"height"`
			ImageURL string `json:"imageUrl"`
			Width    string `json:"width"`
		} `json:"image"`
		ItemAffiliateWebURL string `json:"itemAffiliateWebUrl"`
		ItemCreationDate    string `json:"itemCreationDate"`
		ItemEndDate         string `json:"itemEndDate"`
		ItemGroupHref       string `json:"itemGroupHref"`
		ItemGroupType       string `json:"itemGroupType"`
		ItemHref            string `json:"itemHref"`
		ItemID              string `json:"itemId"`
		ItemLocation        struct {
			AddressLine1    string `json:"addressLine1"`
			AddressLine2    string `json:"addressLine2"`
			City            string `json:"city"`
			Country         string `json:"country"`
			County          string `json:"county"`
			PostalCode      string `json:"postalCode"`
			StateOrProvince string `json:"stateOrProvince"`
		} `json:"itemLocation"`
		ItemWebURL           string   `json:"itemWebUrl"`
		LeafCategoryIds      []string `json:"leafCategoryIds"`
		LegacyItemID         string   `json:"legacyItemId"`
		ListingMarketplaceID string   `json:"listingMarketplaceId"`
		MarketingPrice       struct {
			DiscountAmount struct {
				ConvertedFromCurrency string `json:"convertedFromCurrency"`
				ConvertedFromValue    string `json:"convertedFromValue"`
				Currency              string `json:"currency"`
				Value                 string `json:"value"`
			} `json:"discountAmount"`
			DiscountPercentage string `json:"discountPercentage"`
			OriginalPrice      struct {
				ConvertedFromCurrency string `json:"convertedFromCurrency"`
				ConvertedFromValue    string `json:"convertedFromValue"`
				Currency              string `json:"currency"`
				Value                 string `json:"value"`
			} `json:"originalPrice"`
			PriceTreatment string `json:"priceTreatment"`
		} `json:"marketingPrice"`
		PickupOptions []struct {
			PickupLocationType string `json:"pickupLocationType"`
		} `json:"pickupOptions"`
		Price struct {
			ConvertedFromCurrency string `json:"convertedFromCurrency"`
			ConvertedFromValue    string `json:"convertedFromValue"`
			Currency              string `json:"currency"`
			Value                 string `json:"value"`
		} `json:"price"`
		PriceDisplayCondition string   `json:"priceDisplayCondition"`
		PriorityListing       string   `json:"priorityListing"`
		QualifiedPrograms     []string `json:"qualifiedPrograms"`
		Seller                struct {
			FeedbackPercentage string `json:"feedbackPercentage"`
			FeedbackScore      string `json:"feedbackScore"`
			SellerAccountType  string `json:"sellerAccountType"`
			Username           string `json:"username"`
		} `json:"seller"`
		ShippingOptions []struct {
			GuaranteedDelivery       string `json:"guaranteedDelivery"`
			MaxEstimatedDeliveryDate string `json:"maxEstimatedDeliveryDate"`
			MinEstimatedDeliveryDate string `json:"minEstimatedDeliveryDate"`
			ShippingCost             struct {
				ConvertedFromCurrency string `json:"convertedFromCurrency"`
				ConvertedFromValue    string `json:"convertedFromValue"`
				Currency              string `json:"currency"`
				Value                 string `json:"value"`
			} `json:"shippingCost"`
			ShippingCostType string `json:"shippingCostType"`
		} `json:"shippingOptions"`
		ShortDescription string `json:"shortDescription"`
		ThumbnailImages  []struct {
			Height   string `json:"height"`
			ImageURL string `json:"imageUrl"`
			Width    string `json:"width"`
		} `json:"thumbnailImages"`
		Title                    string `json:"title"`
		TopRatedBuyingExperience string `json:"topRatedBuyingExperience"`
		TyreLabelImageURL        string `json:"tyreLabelImageUrl"`
		UnitPrice                struct {
			ConvertedFromCurrency string `json:"convertedFromCurrency"`
			ConvertedFromValue    string `json:"convertedFromValue"`
			Currency              string `json:"currency"`
			Value                 string `json:"value"`
		} `json:"unitPrice"`
		UnitPricingMeasure string `json:"unitPricingMeasure"`
		WatchCount         string `json:"watchCount"`
	} `json:"itemSummaries"`
	Limit      string `json:"limit"`
	Next       string `json:"next"`
	Offset     string `json:"offset"`
	Prev       string `json:"prev"`
	Refinement struct {
		AspectDistributions []struct {
			AspectValueDistributions []struct {
				LocalizedAspectValue string `json:"localizedAspectValue"`
				MatchCount           string `json:"matchCount"`
				RefinementHref       string `json:"refinementHref"`
			} `json:"aspectValueDistributions"`
			LocalizedAspectName string `json:"localizedAspectName"`
		} `json:"aspectDistributions"`
		BuyingOptionDistributions []struct {
			BuyingOption   string `json:"buyingOption"`
			MatchCount     string `json:"matchCount"`
			RefinementHref string `json:"refinementHref"`
		} `json:"buyingOptionDistributions"`
		CategoryDistributions []struct {
			CategoryID     string `json:"categoryId"`
			CategoryName   string `json:"categoryName"`
			MatchCount     string `json:"matchCount"`
			RefinementHref string `json:"refinementHref"`
		} `json:"categoryDistributions"`
		ConditionDistributions []struct {
			Condition      string `json:"condition"`
			ConditionID    string `json:"conditionId"`
			MatchCount     string `json:"matchCount"`
			RefinementHref string `json:"refinementHref"`
		} `json:"conditionDistributions"`
		DominantCategoryID string `json:"dominantCategoryId"`
	} `json:"refinement"`
	Total    string `json:"total"`
	Warnings []struct {
		Category     string   `json:"category"`
		Domain       string   `json:"domain"`
		ErrorID      string   `json:"errorId"`
		InputRefIds  []string `json:"inputRefIds"`
		LongMessage  string   `json:"longMessage"`
		Message      string   `json:"message"`
		OutputRefIds []string `json:"outputRefIds"`
		Parameters   []struct {
			Name  string `json:"name"`
			Value string `json:"value"`
		} `json:"parameters"`
		Subdomain string `json:"subdomain"`
	} `json:"warnings"`
}
