// Mock the API module before importing the store
const mockLogin = jest.fn();
const mockRegister = jest.fn();
const mockGetProfile = jest.fn();
const mockUpdateProfile = jest.fn();

jest.mock("@/lib/api", () => ({
  authAPI: {
    login: (...args: unknown[]) => mockLogin(...args),
    register: (...args: unknown[]) => mockRegister(...args),
  },
  profileAPI: {
    getProfile: (...args: unknown[]) => mockGetProfile(...args),
    updateProfile: (...args: unknown[]) => mockUpdateProfile(...args),
  },
}));

// Import after mocking
import { useAuthStore } from "../authStore";

// Reset store state between tests
beforeEach(() => {
  useAuthStore.setState({
    token: null,
    user: null,
    profile: null,
    isLoading: false,
    error: null,
  });
  jest.clearAllMocks();
});

describe("authStore - initial state", () => {
  it("has null token initially", () => {
    expect(useAuthStore.getState().token).toBeNull();
  });

  it("has null user initially", () => {
    expect(useAuthStore.getState().user).toBeNull();
  });

  it("has null profile initially", () => {
    expect(useAuthStore.getState().profile).toBeNull();
  });

  it("has isLoading=false initially", () => {
    expect(useAuthStore.getState().isLoading).toBe(false);
  });

  it("has null error initially", () => {
    expect(useAuthStore.getState().error).toBeNull();
  });
});

describe("authStore - login", () => {
  it("sets token and user on successful login", async () => {
    const mockUser = { id: "1", email: "test@example.com", full_name: "Test User" };
    const mockToken = "mock-jwt-token";
    mockLogin.mockResolvedValueOnce({ token: mockToken, user: mockUser });

    await useAuthStore.getState().login("test@example.com", "password123");

    expect(useAuthStore.getState().token).toBe(mockToken);
    expect(useAuthStore.getState().user).toEqual(mockUser);
    expect(useAuthStore.getState().isLoading).toBe(false);
    expect(useAuthStore.getState().error).toBeNull();
  });

  it("sets error on failed login", async () => {
    mockLogin.mockRejectedValueOnce(new Error("Invalid credentials"));

    try {
      await useAuthStore.getState().login("test@example.com", "wrongpassword");
    } catch {
      // Expected to throw
    }

    expect(useAuthStore.getState().error).toBe("Invalid credentials");
    expect(useAuthStore.getState().isLoading).toBe(false);
    expect(useAuthStore.getState().token).toBeNull();
  });

  it("throws error on failed login", async () => {
    mockLogin.mockRejectedValueOnce(new Error("Login failed"));

    await expect(
      useAuthStore.getState().login("test@example.com", "wrongpassword")
    ).rejects.toThrow("Login failed");
  });

  it("calls authAPI.login with correct credentials", async () => {
    mockLogin.mockResolvedValueOnce({ token: "token", user: { id: "1", email: "test@example.com", full_name: "Test" } });

    await useAuthStore.getState().login("test@example.com", "password123");

    expect(mockLogin).toHaveBeenCalledWith("test@example.com", "password123");
  });

  it("sets isLoading=false after successful login", async () => {
    mockLogin.mockResolvedValueOnce({ token: "token", user: { id: "1", email: "test@example.com", full_name: "Test" } });

    await useAuthStore.getState().login("test@example.com", "password123");

    expect(useAuthStore.getState().isLoading).toBe(false);
  });
});

describe("authStore - register", () => {
  it("sets token and user on successful registration", async () => {
    const mockUser = { id: "2", email: "new@example.com", full_name: "New User" };
    const mockToken = "new-jwt-token";
    mockRegister.mockResolvedValueOnce({ token: mockToken, user: mockUser });

    await useAuthStore.getState().register("new@example.com", "password123", "New User");

    expect(useAuthStore.getState().token).toBe(mockToken);
    expect(useAuthStore.getState().user).toEqual(mockUser);
    expect(useAuthStore.getState().isLoading).toBe(false);
  });

  it("sets error on failed registration", async () => {
    mockRegister.mockRejectedValueOnce(new Error("Email already taken"));

    try {
      await useAuthStore.getState().register("existing@example.com", "password123", "User");
    } catch {
      // Expected to throw
    }

    expect(useAuthStore.getState().error).toBe("Email already taken");
    expect(useAuthStore.getState().isLoading).toBe(false);
  });

  it("calls authAPI.register with correct parameters", async () => {
    mockRegister.mockResolvedValueOnce({ token: "token", user: { id: "1", email: "new@example.com", full_name: "New User" } });

    await useAuthStore.getState().register("new@example.com", "password123", "New User");

    expect(mockRegister).toHaveBeenCalledWith("new@example.com", "password123", "New User");
  });
});

describe("authStore - logout", () => {
  it("clears token, user, profile, and error on logout", () => {
    // Set some state first
    useAuthStore.setState({
      token: "some-token",
      user: { id: "1", email: "test@example.com", full_name: "Test User" },
      profile: { id: "1", user_id: "1", headline: "Engineer" } as never,
      error: "some error",
    });

    useAuthStore.getState().logout();

    expect(useAuthStore.getState().token).toBeNull();
    expect(useAuthStore.getState().user).toBeNull();
    expect(useAuthStore.getState().profile).toBeNull();
    expect(useAuthStore.getState().error).toBeNull();
  });
});

describe("authStore - loadProfile", () => {
  it("does nothing when no token", async () => {
    await useAuthStore.getState().loadProfile();

    expect(mockGetProfile).not.toHaveBeenCalled();
  });

  it("loads profile when token is present", async () => {
    const mockProfile = { id: "1", user_id: "1", headline: "Software Engineer" };
    mockGetProfile.mockResolvedValueOnce(mockProfile);
    useAuthStore.setState({ token: "valid-token" });

    await useAuthStore.getState().loadProfile();

    expect(mockGetProfile).toHaveBeenCalledWith("valid-token");
    expect(useAuthStore.getState().profile).toEqual(mockProfile);
    expect(useAuthStore.getState().isLoading).toBe(false);
  });

  it("handles profile load failure gracefully", async () => {
    mockGetProfile.mockRejectedValueOnce(new Error("Network error"));
    useAuthStore.setState({ token: "valid-token" });

    await useAuthStore.getState().loadProfile();

    expect(useAuthStore.getState().isLoading).toBe(false);
    expect(useAuthStore.getState().profile).toBeNull();
  });
});

describe("authStore - updateProfile", () => {
  it("does nothing when no token", async () => {
    await useAuthStore.getState().updateProfile({ headline: "New Headline" });

    expect(mockUpdateProfile).not.toHaveBeenCalled();
  });

  it("updates profile when token is present", async () => {
    const updatedProfile = { id: "1", user_id: "1", headline: "Senior Engineer" };
    mockUpdateProfile.mockResolvedValueOnce(updatedProfile);
    useAuthStore.setState({ token: "valid-token" });

    await useAuthStore.getState().updateProfile({ headline: "Senior Engineer" });

    expect(mockUpdateProfile).toHaveBeenCalledWith("valid-token", { headline: "Senior Engineer" });
    expect(useAuthStore.getState().profile).toEqual(updatedProfile);
  });

  it("sets error on update failure", async () => {
    mockUpdateProfile.mockRejectedValueOnce(new Error("Update failed"));
    useAuthStore.setState({ token: "valid-token" });

    try {
      await useAuthStore.getState().updateProfile({ headline: "New" });
    } catch {
      // Expected to throw
    }

    expect(useAuthStore.getState().error).toBe("Update failed");
    expect(useAuthStore.getState().isLoading).toBe(false);
  });
});

describe("authStore - clearError", () => {
  it("clears the error state", () => {
    useAuthStore.setState({ error: "Some error" });

    useAuthStore.getState().clearError();

    expect(useAuthStore.getState().error).toBeNull();
  });
});
