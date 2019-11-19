package response

func JsonResponse(message string, responseCode int32) map[string]interface{}{
	resp := map[string]interface{}{
		"message": message,
		"statusCode": responseCode,
	}
	return resp
}