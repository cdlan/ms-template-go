package grpc

import (
	"cdlab.cdlan.net/cdlan/uservices/ms-template/internal/config"
	"cdlab.cdlan.net/cdlan/uservices/ms-template/internal/database"
	"cdlab.cdlan.net/cdlan/uservices/ms-template/internal/grpc/gen"
	"context"
	"encoding/json"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"strings"
)

const coverageSqlQuery = `
		SELECT network_operators.name as carrier_name, network_coverage_technologies.name as technology 
		FROM coperture.network_coverage 
		JOIN coperture.network_operators ON network_coverage.network_operator_id = coperture.network_operators.id 
		JOIN coperture.network_coverage_technologies ON network_coverage.network_coverage_technology_id = coperture.network_coverage_technologies.id 
		JOIN coperture.network_coverage_house_numbers ON network_coverage.network_coverage_house_number_id = coperture.network_coverage_house_numbers.id 
		JOIN coperture.network_coverage_addresses ON network_coverage_house_numbers.network_coverage_address_id = coperture.network_coverage_addresses.id 
		JOIN coperture.network_coverage_cities ON network_coverage_addresses.network_coverage_city_id = coperture.network_coverage_cities.id 
		JOIN coperture.network_coverage_states ON network_coverage_cities.network_coverage_state_id = coperture.network_coverage_states.id 
		WHERE coperture.network_coverage_states.abbreviation = $1 AND coperture.network_coverage_cities.name = $2 AND coperture.network_coverage_addresses.name = $3 AND coperture.network_coverage_house_numbers.name = $4
	`

const kitCoverageSQLQuery = `SELECT 
			network_operators.name AS carrier,
    		network_coverage_technologies.name AS technology,
    		commercial_profiles.name,
    		commercial_profiles.upstream_bandwidth,
    		commercial_profiles.downstream_bandwidth
		FROM 
			coperture.network_coverage
		JOIN coperture.network_operators ON network_coverage.network_operator_id = coperture.network_operators.id
		JOIN coperture.network_coverage_technologies ON network_coverage.network_coverage_technology_id = coperture.network_coverage_technologies.id
		JOIN coperture.commercial_profiles ON network_coverage_technologies.id = coperture.commercial_profiles.technology_id
		JOIN coperture.network_coverage_house_numbers ON network_coverage.network_coverage_house_number_id = coperture.network_coverage_house_numbers.id
		JOIN coperture.network_coverage_addresses ON network_coverage_house_numbers.network_coverage_address_id = coperture.network_coverage_addresses.id
		JOIN coperture.network_coverage_cities ON network_coverage_addresses.network_coverage_city_id = coperture.network_coverage_cities.id
		JOIN coperture.network_coverage_states ON network_coverage_cities.network_coverage_state_id = coperture.network_coverage_states.id
		WHERE 
			coperture.network_coverage_states.abbreviation = $1 AND 
			coperture.network_coverage_cities.name = $2 AND 
			coperture.network_coverage_addresses.name = $3 AND 
			coperture.network_coverage_house_numbers.name = $4`

const specificCoverageQuery = `
		SELECT 
			nc.id, no.name as carrier_name, nct.name as technology
		FROM 
			coperture.network_coverage AS nc
		LEFT JOIN coperture.network_operators no ON nc.network_operator_id = no.id
		LEFT JOIN coperture.network_coverage_technologies nct ON nc.network_coverage_technology_id = nct.id
		LEFT JOIN coperture.network_coverage_house_numbers numbers ON nc.network_coverage_house_number_id = numbers.id
		LEFT JOIN coperture.network_coverage_addresses addresses ON numbers.network_coverage_address_id = addresses.id
		LEFT JOIN coperture.network_coverage_cities cities ON addresses.network_coverage_city_id = cities.id
		LEFT JOIN coperture.network_coverage_states provinces ON cities.network_coverage_state_id = provinces.id
		LEFT JOIN coperture.commercial_profiles cp ON nct.id = cp.technology_id
		WHERE 
			provinces.abbreviation = $1 AND 
			cities.name = $2 AND 
			addresses.name = $3 AND 
			numbers.name = $4 AND
			no.name = ANY($5) AND 
			cp.name = ANY($6)
	`

type CoverageServer struct {
	gen.UnimplementedCoverageServiceServer
}

func NewCoverageServerServer() *CoverageServer {
	return &CoverageServer{}
}

func (S *CoverageServer) GetCoverage(req *gen.GetCoverageRequest, stream gen.CoverageService_GetCoverageServer) error {
	// start span
	ctx, span := config.C.Otel.NewSpan(stream.Context(), "main.CoverageServer/GetCoverage")
	if span != nil {
		defer span.End()
	}

	if database.DBState.DBPool == nil {

		err := fmt.Errorf("database pool is nil")
		config.C.Otel.RecordError(span, err)
		return status.Errorf(codes.Internal, "%v", err)
	}

	// Fetching the address from the request
	address := req.GetAddress()

	sqlQuery := coverageSqlQuery

	rows, err := database.DBState.DBPool.Query(ctx, sqlQuery, address.GetProvince(), address.GetCity(), address.GetAddress(), address.GetHouseNumber())
	if err != nil {
		return err
	}
	defer rows.Close()

	var coverageInfos []*gen.CoverageInfo
	for rows.Next() {
		var info gen.CoverageInfo
		err := rows.Scan(&info.CarrierName, &info.Technology)
		if err != nil {
			return err
		}
		coverageInfos = append(coverageInfos, &info)
	}

	if err = rows.Err(); err != nil {
		return err
	}

	return nil
}

func (S *CoverageServer) GetCoverageKit(req *gen.GetCoverageKitRequest, stream gen.CoverageService_GetCoverageKitServer) error {
	// start span
	ctx, span := config.C.Otel.NewSpan(stream.Context(), "main.CoverageServer/GetKitCoverage")
	if span != nil {
		defer span.End()
	}

	if database.DBState.DBPool == nil {

		err := fmt.Errorf("database pool is nil")
		config.C.Otel.RecordError(span, err)
		return status.Errorf(codes.Internal, "%v", err)
	}

	// Fetching the address from the request
	address := req.GetAddress()

	sqlQuery := kitCoverageSQLQuery

	rows, err := database.DBState.DBPool.Query(ctx, sqlQuery, address.GetProvince(), address.GetCity(), address.GetAddress(), address.GetHouseNumber())
	if err != nil {
		return err
	}
	defer rows.Close()

	var kitInfos []*gen.KitInfo

	for rows.Next() {
		var kit gen.KitInfo
		var commercialProfile gen.CommercialProfile

		if err := rows.Scan(
			&kit.Carrier,
			&kit.Technology,
			&commercialProfile.Name,
			&commercialProfile.UpstreamBandwidth,
			&commercialProfile.DownstreamBandwidth,
			// Add other fields from commercial_profiles as needed
		); err != nil {
			return err
		}

		kit.CommercialProfile = &commercialProfile
		kitInfos = append(kitInfos, &kit)
	}

	if err = rows.Err(); err != nil {
		return err
	}
	// Build the response
	response := &gen.GetCoverageKitResponse{
		Kits: kitInfos,
	}
	// Send the response over the stream
	if err = stream.Send(response); err != nil {
		span.RecordError(err)
		return status.Errorf(codes.Internal, "Failed to send response: %v", err)
	}
	return nil
}

func (S *CoverageServer) QuerySpecificCoverage(req *gen.QuerySpecificCoverageRequest, stream gen.CoverageService_QuerySpecificCoverageServer) error {
	// start span
	ctx, span := config.C.Otel.NewSpan(stream.Context(), "mmain.CoverageServer/QuerySpecificCoverage")
	if span != nil {
		defer span.End()
	}

	if database.DBState.DBPool == nil {
		err := fmt.Errorf("database pool is nil")
		config.C.Otel.RecordError(span, err)
		return status.Errorf(codes.Internal, "%v", err)
	}

	// Use the core logic to fetch specific coverage from the DB
	coverageInfos, err := fetchSpecificCoverageFromDB(ctx, req.GetAddress(), req.GetCommercialProfiles(), req.GetAvailableCarriers())
	if err != nil {
		config.C.Otel.RecordError(span, err)
		return status.Errorf(codes.Internal, "Failed to fetch specific coverage from DB: %v", err)
	}

	// Build and send the response over the stream
	response := &gen.QuerySpecificCoverageResponse{
		CoverageInfos: coverageInfos,
	}
	if err := stream.Send(response); err != nil {
		config.C.Otel.RecordError(span, err)
		return status.Errorf(codes.Internal, "Failed to send response: %v", err)
	}

	return nil
}

// Moved the core function here for clarity, you can keep it in your server's methods if you prefer
func fetchSpecificCoverageFromDB(ctx context.Context, address *gen.Address, commercialProfiles []string, availableCarriers []string) ([]*gen.CoverageInfo, error) {
	var coverageInfos []*gen.CoverageInfo
	query := specificCoverageQuery

	rows, err := database.DBState.DBPool.Query(context.Background(), query,
		strings.ToUpper(address.GetProvince()),
		address.GetCity(),
		address.GetAddress(),
		address.GetHouseNumber(),
		availableCarriers,
		commercialProfiles,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var info gen.CoverageInfo
		if err := rows.Scan(&info.CarrierName, &info.Technology); err != nil {
			return nil, err
		}
		coverageInfos = append(coverageInfos, &info)
	}

	return coverageInfos, nil
}

func (S *CoverageServer) GetStates(empty *emptypb.Empty, stream gen.CoverageService_GetStatesServer) error {
	// start span
	ctx, span := config.C.Otel.NewSpan(stream.Context(), "main.CoverageServer/GetStates")
	if span != nil {
		defer span.End()
	}

	if database.DBState.DBPool == nil {
		err := fmt.Errorf("database pool is nil")
		config.C.Otel.RecordError(span, err)
		return status.Errorf(codes.Internal, "%v", err)
	}

	// Use the CALL command to call the stored procedure
	sqlQuery := `SELECT coperture.get_states()`

	// Execute the procedure call
	rows, err := database.DBState.DBPool.Query(ctx, sqlQuery)
	if err != nil {
		config.C.Otel.RecordError(span, err)
		return status.Errorf(codes.Internal, "Failed to execute SQL query: %v", err)
	}
	defer rows.Close()

	// Iterate over the result rows, unmarshal the JSON, and map them to gRPC response structure
	var states []*gen.State
	for rows.Next() {
		var stateData string
		if err := rows.Scan(&stateData); err != nil {
			config.C.Otel.RecordError(span, err)
			return status.Errorf(codes.Internal, "Failed to scan row data: %v", err)
		}

		// Assuming stateData contains a JSON array of states
		var statesSlice []gen.State
		if err := json.Unmarshal([]byte(stateData), &statesSlice); err != nil {
			config.C.Otel.RecordError(span, err)
			return status.Errorf(codes.Internal, "Failed to unmarshal state data: %v", err)
		}
		for _, state := range statesSlice {
			state := state
			states = append(states, &state)
		}
	}

	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		config.C.Otel.RecordError(span, err)
		return status.Errorf(codes.Internal, "Error iterating over rows: %v", err)
	}

	// Build and send the response over the stream
	response := &gen.GetStatesResponse{
		States: states,
	}
	if err := stream.Send(response); err != nil {
		config.C.Otel.RecordError(span, err)
		return status.Errorf(codes.Internal, "Failed to send response: %v", err)
	}

	return nil
}

//func (S *ExampleServer) UnaryMethod(ctx context.Context, req *gen.Get) (*gen.Item, error) {
//
//	// start span
//	ctx, span := config.C.Otel.NewSpan(ctx, "main.Example/UnaryMethod")
//	if span != nil {
//		defer span.End()
//	}
//
//	// do stuff
//
//	return &gen.Item{}, nil
//}
