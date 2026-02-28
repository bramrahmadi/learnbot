import { create } from "zustand";
import { persist } from "zustand/middleware";
import { authAPI, profileAPI, type UserProfile } from "@/lib/api";

interface AuthState {
  token: string | null;
  user: { id: string; email: string; full_name: string } | null;
  profile: UserProfile | null;
  isLoading: boolean;
  error: string | null;
  login: (email: string, password: string) => Promise<void>;
  register: (email: string, password: string, fullName: string) => Promise<void>;
  logout: () => void;
  loadProfile: () => Promise<void>;
  updateProfile: (data: Partial<UserProfile>) => Promise<void>;
  clearError: () => void;
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set, get) => ({
      token: null, user: null, profile: null, isLoading: false, error: null,
      login: async (email, password) => {
        set({ isLoading: true, error: null });
        try {
          const result = await authAPI.login(email, password);
          set({ token: result.token, user: result.user, isLoading: false });
        } catch (err: unknown) {
          set({ isLoading: false, error: err instanceof Error ? err.message : "Login failed" });
          throw err;
        }
      },
      register: async (email, password, fullName) => {
        set({ isLoading: true, error: null });
        try {
          const result = await authAPI.register(email, password, fullName);
          set({ token: result.token, user: result.user, isLoading: false });
        } catch (err: unknown) {
          set({ isLoading: false, error: err instanceof Error ? err.message : "Registration failed" });
          throw err;
        }
      },
      logout: () => set({ token: null, user: null, profile: null, error: null }),
      loadProfile: async () => {
        const { token } = get();
        if (!token) return;
        set({ isLoading: true });
        try {
          const profile = await profileAPI.getProfile(token);
          set({ profile, isLoading: false });
        } catch { set({ isLoading: false }); }
      },
      updateProfile: async (data) => {
        const { token } = get();
        if (!token) return;
        set({ isLoading: true, error: null });
        try {
          const profile = await profileAPI.updateProfile(token, data);
          set({ profile, isLoading: false });
        } catch (err: unknown) {
          set({ isLoading: false, error: err instanceof Error ? err.message : "Update failed" });
          throw err;
        }
      },
      clearError: () => set({ error: null }),
    }),
    { name: "learnbot-auth", partialize: (state) => ({ token: state.token, user: state.user }) }
  )
);

export const useToken = () => useAuthStore((s) => s.token);
export const useUser = () => useAuthStore((s) => s.user);
export const useIsAuthenticated = () => useAuthStore((s) => !!s.token);
