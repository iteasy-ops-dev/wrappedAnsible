package model

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	config "iteasy.wrappedAnsible/configs"
)

type WorkHistory struct {
	Index                   string `bson:"index"`
	Status                  string `bson:"status"`
	RegistrationDate        string `bson:"registration_date"`
	DesiredWorkingHours     string `bson:"desired_working_hours"`
	EstimatedCompletionTime string `bson:"estimated_completion_time"`
	Worker                  string `bson:"worker"`
	WorkRequestItems        string `bson:"work_request_items"`
	WorkRequestDetails      string `bson:"work_request_details"`
	SubCategory             string `bson:"sub_category"`
	ClientCompany           string `bson:"client_company"`
	IP                      string `bson:"ip"`
	Brand                   string `bson:"brand"`
	Section                 string `bson:"section"`
	Url                     string `bson:"url"`

	baseUrl string
}

func (w *WorkHistory) _collection() *mongo.Collection {
	return db.Collection(config.GlobalConfig.MongoDB.Collections.WorkHistory)
}

func NewWorkHistory() *WorkHistory {
	return &WorkHistory{
		baseUrl: config.GlobalConfig.Erp.BaseURL,
	}
}

func (w *WorkHistory) Setindex(index string) {
	w.Index = index
}
func (w *WorkHistory) Setstatus(status string) {
	w.Status = status
}
func (w *WorkHistory) SetregistrationDate(registrationDate string) {
	w.RegistrationDate = registrationDate
}
func (w *WorkHistory) SetdesiredWorkingHours(desiredWorkingHours string) {
	w.DesiredWorkingHours = desiredWorkingHours
}
func (w *WorkHistory) SetestimatedCompletionTime(estimatedCompletionTime string) {
	w.EstimatedCompletionTime = estimatedCompletionTime
}
func (w *WorkHistory) Setworker(worker string) {
	w.Worker = worker
}
func (w *WorkHistory) SetworkRequestItems(workRequestItems string) {
	w.WorkRequestItems = workRequestItems
}
func (w *WorkHistory) SetworkRequestDetails(workRequestDetails string) {
	w.WorkRequestDetails = workRequestDetails
}
func (w *WorkHistory) SetsubCategory(subCategory string) {
	w.SubCategory = subCategory
}
func (w *WorkHistory) SetclientCompany(clientCompany string) {
	w.ClientCompany = clientCompany
}
func (w *WorkHistory) Setip(ip string) {
	w.IP = ip
}
func (w *WorkHistory) Setbrand(brand string) {
	w.Brand = brand
}
func (w *WorkHistory) Setsection(section string) {
	w.Section = section
}
func (w *WorkHistory) Seturl(url string) {
	w.Url = fmt.Sprintf("%s%s", w.baseUrl, url)
}

func (w *WorkHistory) Print() {
	fmt.Println(w)
}

// 인덱스가 숫자로만 구성된 문자열인지 확인
func (w *WorkHistory) isIndexValid() bool {
	_, err := strconv.Atoi(w.Index)
	return err == nil
}

func (w *WorkHistory) isValid() bool {
	valid := w.isIndexValid() &&
		strings.TrimSpace(w.Status) != "" &&
		strings.TrimSpace(w.RegistrationDate) != "" &&
		// strings.TrimSpace(w.EstimatedCompletionTime) != "" &&
		strings.TrimSpace(w.WorkRequestItems) != "" &&
		// strings.TrimSpace(w.SubCategory) != "" &&
		strings.TrimSpace(w.ClientCompany) != "" &&
		strings.TrimSpace(w.Brand) != "" &&
		strings.TrimSpace(w.Section) != "" &&
		strings.TrimSpace(w.Url) != ""

	if !valid {
		w.PrintValidations() // 유효성 검사 실패 시 필드 값 출력
	}
	return valid
}

func (w *WorkHistory) PrintValidations() {
	fmt.Println("==================  PrintValidations  =========================")
	fmt.Printf("Index: '%s' Valid: %v\n", w.Index, w.isIndexValid())
	fmt.Printf("Status: '%s' Valid: %v\n", w.Status, strings.TrimSpace(w.Status) != "")
	fmt.Printf("RegistrationDate: '%s' Valid: %v\n", w.RegistrationDate, strings.TrimSpace(w.RegistrationDate) != "")
	// fmt.Printf("EstimatedCompletionTime: '%s' Valid: %v\n", w.EstimatedCompletionTime, strings.TrimSpace(w.EstimatedCompletionTime) != "")
	fmt.Printf("WorkRequestItems: '%s' Valid: %v\n", w.WorkRequestItems, strings.TrimSpace(w.WorkRequestItems) != "")
	// fmt.Printf("SubCategory: '%s' Valid: %v\n", w.SubCategory, strings.TrimSpace(w.SubCategory) != "")
	fmt.Printf("ClientCompany: '%s' Valid: %v\n", w.ClientCompany, strings.TrimSpace(w.ClientCompany) != "")
	fmt.Printf("Brand: '%s' Valid: %v\n", w.Brand, strings.TrimSpace(w.Brand) != "")
	fmt.Printf("Section: '%s' Valid: %v\n", w.Section, strings.TrimSpace(w.Section) != "")
	fmt.Printf("URL: '%s' Valid: %v\n", w.Url, strings.TrimSpace(w.Url) != "")
}

func (w *WorkHistory) Put() error {
	if !w.isValid() {
		return fmt.Errorf("no input")
	}

	col := w._collection()

	// 필드 일치 여부를 확인하기 위한 필터를 작성
	filter := bson.M{
		// "index":                     w.Index, // 인덱스가 변할 수도 있었음.
		// "registration_date":         w.RegistrationDate,
		// "client_company":            w.ClientCompany,
		// "sub_category":              w.SubCategory,
		// "work_request_items":        w.WorkRequestItems,
		// "estimated_completion_time": w.EstimatedCompletionTime,
		"url": w.Url,
	}

	// 해당 필터로 이미 존재하는지 확인
	existingReq := &WorkHistory{}
	err := col.FindOne(context.Background(), filter).Decode(existingReq)
	if err == nil {
		// 일부 데이터만 일치하면 상태를 업데이트
		update := bson.M{
			"$set": bson.M{
				// "index":   w.Index,
				// "status":  w.Status,
				// "worker":  w.Worker,
				// "ip":      w.IP,
				// "brand":   w.Brand,
				// "section": w.Section,
				"index":                     w.Index,
				"status":                    w.Status,
				"registration_date":         w.RegistrationDate,
				"desired_working_hours":     w.DesiredWorkingHours,
				"estimated_completion_time": w.EstimatedCompletionTime,
				"worker":                    w.Worker,
				"work_request_items":        w.WorkRequestItems,
				"work_request_details":      w.WorkRequestDetails,
				"sub_category":              w.SubCategory,
				"client_company":            w.ClientCompany,
				"ip":                        w.IP,
				"brand":                     w.Brand,
				"section":                   w.Section,
				// "url":     w.Url,
			},
		}
		_, err = col.UpdateOne(context.Background(), filter, update)
		// fmt.Printf("업데이트완료: %s\n", w.Url)
		return err
	}
	// else {
	// 	fmt.Printf("필터 검색 에러: %s\n%s\n", w.Url, err)
	// }

	if err == mongo.ErrNoDocuments {
		// 완전히 일치하는 데이터가 없으면 새 데이터 삽입
		_, err = col.InsertOne(context.Background(), w)
		// fmt.Printf("삽입완료: %s\n", w.Url)
		return err
	}

	return err
}

func (w *WorkHistory) Get(filter bson.M, page int, pageSize int) ([]WorkHistory, int, error) {
	col := db.Collection(config.GlobalConfig.MongoDB.Collections.WorkHistory)

	// 페이지 번호와 페이지 크기 검증
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10 // 기본값 설정
	}

	// Skip 및 Limit 설정
	skip := (page - 1) * pageSize

	// 정렬 옵션 추가
	sortOptions := bson.D{
		{Key: "registration_date", Value: -1}, // RegistrationDate 내림차순
		{Key: "index", Value: -1},             // Index 내림차순
	}

	// 데이터 조회
	cursor, err := col.Find(
		context.Background(),
		filter,
		options.Find().
			SetSkip(int64(skip)).
			SetLimit(int64(pageSize)).
			SetSort(sortOptions), // 정렬 옵션 설정
	)
	if err != nil {
		return nil, 0, err
	}

	// 커서를 디코딩하여 결과 반환
	results, err := DecodeCursor[WorkHistory](cursor)
	if err != nil {
		return nil, 0, err
	}

	// 전체 데이터 수 조회
	totalCount, err := col.CountDocuments(context.Background(), filter)
	if err != nil {
		return nil, 0, err
	}

	totalPages := int((totalCount + int64(pageSize) - 1) / int64(pageSize)) // 총 페이지 수 계산

	return results, totalPages, nil
}

func (w *WorkHistory) Dashboard(filter bson.M) (map[string]interface{}, error) {
	col := w._collection()

	// 통계 결과를 담을 map 초기화
	stats := make(map[string]interface{})

	// 1. WorkRequestItems, SubCategory, ClientCompany, IP 별 통계 조회
	columns := []string{"work_request_items", "sub_category", "client_company", "ip"}

	for _, column := range columns {
		pipeline := []bson.M{
			{"$match": filter}, // 필터 적용 (필요 시 조건 추가)
			{"$group": bson.M{
				"_id":   fmt.Sprintf("$%s", column), // 해당 컬럼의 각 값을 그룹화
				"count": bson.M{"$sum": 1},          // 각 그룹의 개수를 셈
			}},
			{"$sort": bson.M{"count": -1}}, // 개수 순으로 정렬
		}

		cursor, err := col.Aggregate(context.Background(), pipeline)
		if err != nil {
			return nil, err
		}

		var results []bson.M
		if err = cursor.All(context.Background(), &results); err != nil {
			return nil, err
		}

		// 총 개수 구하기
		totalCount := 0
		for _, result := range results {
			switch v := result["count"].(type) {
			case int32:
				totalCount += int(v)
			case int64:
				totalCount += int(v)
			case int:
				totalCount += v
			default:
				// 처리할 수 없는 타입의 경우 에러 처리
				return nil, fmt.Errorf("unexpected count type: %T", v)
			}
		}

		// 컬럼별 통계 및 총 개수 추가
		stats[column] = bson.M{
			"total": totalCount,
			"data":  results,
		}
	}

	// 2. 데이터 기간 조회 (RegistrationDate 기준)
	datePipeline := []bson.M{
		{"$match": filter}, // 필터 적용 (없으면 모든 데이터 대상)
		{"$group": bson.M{
			"_id": nil,
			"minDate": bson.M{
				"$min": "$registration_date",
			},
			"maxDate": bson.M{
				"$max": "$registration_date",
			},
		}},
	}

	dateCursor, err := col.Aggregate(context.Background(), datePipeline)
	if err != nil {
		return nil, err
	}

	type DateResult struct {
		MinDate string `bson:"minDate"`
		MaxDate string `bson:"maxDate"`
	}

	var dateResult DateResult
	if dateCursor.Next(context.Background()) {
		if err := dateCursor.Decode(&dateResult); err != nil {
			return nil, err
		}
		stats["DataPeriod"] = bson.M{
			"From": dateResult.MinDate,
			"To":   dateResult.MaxDate,
		}
	}

	return stats, nil
}
