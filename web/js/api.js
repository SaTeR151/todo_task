const ACCESS_TOKEN_KEY = "todo_task.access_token";
const REFRESH_TOKEN_KEY = "todo_task.refresh_token";

export class ApiError extends Error {
  constructor(message, status, payload) {
    super(message);
    this.name = "ApiError";
    this.status = status;
    this.payload = payload;
  }
}

export function hasToken() {
  return Boolean(localStorage.getItem(ACCESS_TOKEN_KEY));
}

export function saveSession(tokens) {
  localStorage.setItem(ACCESS_TOKEN_KEY, tokens.access_token);
  localStorage.setItem(REFRESH_TOKEN_KEY, tokens.refresh_token);
}

export function clearSession() {
  localStorage.removeItem(ACCESS_TOKEN_KEY);
  localStorage.removeItem(REFRESH_TOKEN_KEY);
}

function getAccessToken() {
  return localStorage.getItem(ACCESS_TOKEN_KEY);
}

function getRefreshToken() {
  return localStorage.getItem(REFRESH_TOKEN_KEY);
}

function buildMessage(payload, fallback) {
  if (typeof payload === "string" && payload.trim()) {
    return payload;
  }

  if (payload && typeof payload.error === "string") {
    return payload.error;
  }

  if (payload && typeof payload.debug === "string") {
    return payload.debug;
  }

  return fallback;
}

async function parseResponse(response) {
  if (response.status === 204) {
    return null;
  }

  const text = await response.text();
  if (!text) {
    return null;
  }

  try {
    return JSON.parse(text);
  } catch {
    return text;
  }
}

async function request(path, options = {}) {
  const {
    method = "GET",
    body,
    auth = true,
    emptyOnNotFound = false,
  } = options;

  const headers = {};
  const token = getAccessToken();

  if (body !== undefined) {
    headers["Content-Type"] = "application/json";
  }

  if (auth && token) {
    headers.Authorization = `Bearer ${token}`;
  }

  const response = await fetch(`/api${path}`, {
    method,
    headers,
    body: body === undefined ? undefined : JSON.stringify(body),
  });

  const payload = await parseResponse(response);

  if (response.status === 404 && emptyOnNotFound) {
    return [];
  }

  if (response.status === 401 && auth) {
    clearSession();
  }

  if (!response.ok) {
    throw new ApiError(
      buildMessage(payload, `HTTP ${response.status}`),
      response.status,
      payload,
    );
  }

  return payload;
}

function cleanPatch(payload) {
  return Object.fromEntries(
    Object.entries(payload).filter(([, value]) => value !== undefined),
  );
}

export const api = {
  auth: {
    login(login, password) {
      return request("/user/auth", {
        method: "POST",
        auth: false,
        body: { login, password },
      });
    },
    register(login, password) {
      return request("/user/", {
        method: "POST",
        auth: false,
        body: { login, password },
      });
    },
    currentUser() {
      return request("/user/");
    },
    changePassword(oldPassword, newPassword) {
      return request("/user/password-change", {
        method: "PATCH",
        body: {
          old_password: oldPassword,
          new_password: newPassword,
        },
      });
    },
    async logout() {
      try {
        await request("/user/logout", { method: "POST" });
      } finally {
        clearSession();
      }
    },
    refresh() {
      return request("/user/refresh", {
        method: "POST",
        body: { refresh_token: getRefreshToken() },
      });
    },
  },
  boards: {
    list() {
      return request("/boards/", { emptyOnNotFound: true });
    },
    get(boardID) {
      return request(`/boards/${encodeURIComponent(boardID)}`);
    },
    create(name) {
      return request("/boards/", {
        method: "POST",
        body: { name },
      });
    },
    update(boardID, payload) {
      return request(`/boards/${encodeURIComponent(boardID)}`, {
        method: "PATCH",
        body: cleanPatch(payload),
      });
    },
    delete(boardID) {
      return request(`/boards/${encodeURIComponent(boardID)}`, {
        method: "DELETE",
      });
    },
  },
  types: {
    list() {
      return request("/types/", { emptyOnNotFound: true });
    },
    create(name, color) {
      return request("/types/", {
        method: "POST",
        body: { name, color },
      });
    },
    update(typeID, payload) {
      return request(`/types/${encodeURIComponent(typeID)}`, {
        method: "PATCH",
        body: cleanPatch(payload),
      });
    },
    delete(typeID) {
      return request(`/types/${encodeURIComponent(typeID)}`, {
        method: "DELETE",
      });
    },
  },
  columns: {
    list(boardID) {
      return request(`/boards/${encodeURIComponent(boardID)}/columns/`, {
        emptyOnNotFound: true,
      });
    },
    create(boardID, payload) {
      return request(`/boards/${encodeURIComponent(boardID)}/columns/`, {
        method: "POST",
        body: payload,
      });
    },
    update(boardID, columnID, payload) {
      return request(
        `/boards/${encodeURIComponent(boardID)}/columns/${encodeURIComponent(columnID)}`,
        {
          method: "PATCH",
          body: cleanPatch(payload),
        },
      );
    },
    delete(boardID, columnID) {
      return request(
        `/boards/${encodeURIComponent(boardID)}/columns/${encodeURIComponent(columnID)}`,
        { method: "DELETE" },
      );
    },
    swap(boardID, columnA, columnB) {
      return request(`/boards/${encodeURIComponent(boardID)}/columns/swap`, {
        method: "PUT",
        body: {
          column_a: columnA,
          column_b: columnB,
        },
      });
    },
  },
  tasks: {
    list(boardID, filters = {}) {
      const params = new URLSearchParams();
      if (filters.column_id) {
        params.set("column_id", filters.column_id);
      }
      if (filters.type_id) {
        params.set("type_id", filters.type_id);
      }
      const query = params.toString();
      return request(
        `/boards/${encodeURIComponent(boardID)}/tasks/${query ? `?${query}` : ""}`,
        { emptyOnNotFound: true },
      );
    },
    get(boardID, taskID) {
      return request(
        `/boards/${encodeURIComponent(boardID)}/tasks/${encodeURIComponent(taskID)}`,
      );
    },
    create(boardID, payload) {
      return request(`/boards/${encodeURIComponent(boardID)}/tasks/`, {
        method: "POST",
        body: payload,
      });
    },
    update(boardID, taskID, payload) {
      return request(
        `/boards/${encodeURIComponent(boardID)}/tasks/${encodeURIComponent(taskID)}`,
        {
          method: "PATCH",
          body: cleanPatch(payload),
        },
      );
    },
    delete(boardID, taskID) {
      return request(
        `/boards/${encodeURIComponent(boardID)}/tasks/${encodeURIComponent(taskID)}`,
        { method: "DELETE" },
      );
    },
    move(boardID, taskID, columnID) {
      return request(
        `/boards/${encodeURIComponent(boardID)}/tasks/${encodeURIComponent(taskID)}/move`,
        {
          method: "PUT",
          body: { column_id: columnID },
        },
      );
    },
  },
};
