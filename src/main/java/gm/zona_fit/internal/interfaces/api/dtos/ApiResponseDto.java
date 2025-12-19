package gm.zona_fit.internal.interfaces.api.dtos;

import com.fasterxml.jackson.annotation.JsonProperty;

public class ApiResponseDto<T> {
    
    @JsonProperty("success")
    private boolean success;
    
    @JsonProperty("code")
    private int code;
    
    @JsonProperty("data")
    private T data;
    
    @JsonProperty("message")
    private String message;
    
    // Constructores
    public ApiResponseDto(boolean success, int code, T data) {
        this.success = success;
        this.code = code;
        this.data = data;
    }
    
    public ApiResponseDto(boolean success, int code, String message) {
        this.success = success;
        this.code = code;
        this.message = message;
    }
    
    public ApiResponseDto(boolean success, int code, T data, String message) {
        this.success = success;
        this.code = code;
        this.data = data;
        this.message = message;
    }
    
    // Métodos estáticos de conveniencia
    public static <T> ApiResponseDto<T> success(T data) {
        return new ApiResponseDto<>(true, 200, data);
    }
    
    public static <T> ApiResponseDto<T> success(T data, String message) {
        return new ApiResponseDto<>(true, 200, data, message);
    }
    
    public static <T> ApiResponseDto<T> error(int code, String message) {
        return new ApiResponseDto<>(false, code, message);
    }
    
    // Getters
    public boolean isSuccess() { return success; }
    public int getCode() { return code; }
    public T getData() { return data; }
    public String getMessage() { return message; }
}