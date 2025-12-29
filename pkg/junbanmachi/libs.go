package junbanmachi // 順番待ち - じゅんばんまち - Queuing

func pushRequest2Packet(pushRequest *PushRequest) *Packet {
	retPacket := &Packet{
		ID:       pushRequest.ID,
		SourceIP: pushRequest.SourceIP,
		Domain:   pushRequest.Domain,
		URI:      pushRequest.URI,
		Method:   pushRequest.Method,
		//Service:  pushRequest.Domain, // thêm service để phân loại queue
		Headers:  make([]*PushRequest_HeaderSchema, 0),
	}

	for _, header := range pushRequest.Headers {
		retPacket.Headers = append(retPacket.Headers, &PushRequest_HeaderSchema{
			Field: header.Field,
			Value: header.Value,
		})
	}

	return retPacket
}

func packet2PopResponse(packet *Packet) *PopResponse {
	retPopResponse := &PopResponse{
		ID:       packet.ID,
		SourceIP: packet.SourceIP,
		Domain:   packet.Domain,
		URI:      packet.URI,
		Method:   packet.Method,
		Headers:  make([]*PopResponse_HeaderSchema, 0),
	}

	for _, header := range packet.Headers {
		retPopResponse.Headers = append(retPopResponse.Headers, &PopResponse_HeaderSchema{
			Field: header.Field,
			Value: header.Value,
		})
	}

	return retPopResponse
}
