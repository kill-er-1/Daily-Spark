const BASE_URL =
  import.meta.env.VITE_API_BASE_URL ?? "http://localhost:8080/api/v1";

interface ApiResponse<T> {
  data?: T;//？表示可选
  error?: string;
  message?: string;
}

async function request<T = any>(
  path: string,
  options: RequestInit = {}
): Promise<T> {
  const res = await fetch(`${BASE_URL}${path}`, {
    headers: {
      "Content-Type": "application/json",
      ...(options.headers || {}),
    },
    ...options,
  });

  let data: ApiResponse<T> | null = null;
  try {
    data = await res.json();
  } catch {
    data = null;
  }

  if (!res.ok) {
    const errorMsg = data?.error || data?.message || `HTTP ${res.status}`;
    throw new Error(errorMsg);
  }

  return (data as T) ?? ({} as T);
}

interface RegisterPayload {
  account: string;
  password: string;
}

interface LoginPayload {
  account: string;
  password: string;
}

interface UpdatePasswordPayload {
  password: string;
}

interface DeleteUserPayload {
  account: string;
}

export const api = {
  register(payload: RegisterPayload) {
    return request("/users/signup", {
      method: "POST",
      body: JSON.stringify(payload),
    });
  },

  login(payload: LoginPayload) {
    return request("/users/signin", {
      method: "POST",
      body: JSON.stringify(payload),
    });
  },

  updatePassword(userId: string, payload: UpdatePasswordPayload) {
    return request(`/users/update/${userId}`, {
      method: "POST",
      body: JSON.stringify(payload),
    });
  },

  deleteUser(userId: string, payload: DeleteUserPayload) {
    return request(`/users/delete/${userId}`, {
      method: "POST",
      body: JSON.stringify(payload),
    });
  },
};
