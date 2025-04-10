import { useLocation, useNavigate } from '@solidjs/router';
import { Component, createSignal, onMount } from 'solid-js';
import { loginApi } from '../api/login';
import { userApi } from '../api/user';
import { Input } from '@/components/ui/input';

const CreateLogin: Component = () => {
  const navigate = useNavigate();
  const location = useLocation();
  const [username, setUsername] = createSignal('');
  const [password, setPassword] = createSignal('');
  const [confirmPassword, setConfirmPassword] = createSignal('');
  const [error, setError] = createSignal<string | null>(null);
  const [loading, setLoading] = createSignal(false);
  const [userId, setUserId] = createSignal<string | null>(null);
  const [returnToUsers, setReturnToUsers] = createSignal(false);
  const [userData, setUserData] = createSignal<any>(null);

  onMount(() => {
    // Check if we have a userId in the query parameters
    const params = new URLSearchParams(location.search);
    const userIdParam = params.get('userId');
    if (userIdParam) {
      setUserId(userIdParam);
      setReturnToUsers(true);
      
      // Fetch user data to pre-fill fields
      fetchUserData(userIdParam);
    }
  });
  
  const fetchUserData = async (userId: string) => {
    try {
      const userData = await userApi.getUser(userId);
      if (userData) {
        // Store the user data
        setUserData(userData);
        // Pre-fill fields with user data
        setUsername(userData.username || '');
      }
    } catch (err) {
      console.error('Failed to fetch user data:', err);
    }
  };

  const handleSubmit = async (e: Event) => {
    e.preventDefault();
    
    // Validate form
    if (password() !== confirmPassword()) {
      setError('Passwords do not match');
      return;
    }

    setError(null);
    setLoading(true);

    try {
      const newLogin = await loginApi.createLogin({
        username: username(),
        password: password(),
      });
      
      // If we have a userId, associate the login with the user
      if (userId() && newLogin.id) {
        const user = userData();
        await userApi.updateUser(userId()!, {
          login_id: newLogin.id,
          name: user?.name,
          username: user?.username,
          role_ids: user?.roles?.map((role: any) => role.id) || []
        });
      }
      
      // Navigate back to the appropriate page
      if (returnToUsers()) {
        navigate('/users');
      } else {
        navigate('/logins');
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to create login');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div class="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
      <div class="md:flex md:items-center md:justify-between">
        <div class="min-w-0 flex-1">
          <h2 class="text-2xl font-bold leading-7 text-gray-12 sm:truncate sm:text-3xl sm:tracking-tight">
            {userId() ? 'Create Login for User' : 'Create New Login'}
          </h2>
        </div>
        <div class="mt-4 flex md:ml-4 md:mt-0">
          <button
            type="button"
            onClick={() => navigate(-1)}
            class="inline-flex items-center rounded-lg bg-white px-3 py-2 text-sm font-semibold text-gray-11 shadow-sm ring-1 ring-inset ring-gray-6 hover:bg-gray-3"
          >
            Cancel
          </button>
        </div>
      </div>

      <div class="mt-8 flow-root">
        <div class="overflow-hidden bg-white shadow rounded-lg">
          <form onSubmit={handleSubmit} class="px-4 py-5 sm:p-6">
            {error() && (
              <div class="mb-4 rounded-lg bg-red-50 p-4">
                <div class="flex">
                  <div class="flex-shrink-0">
                    <svg class="h-5 w-5 text-red-400" viewBox="0 0 20 20" fill="currentColor">
                      <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd" />
                    </svg>
                  </div>
                  <div class="ml-3">
                    <h3 class="text-sm font-medium text-red-800">{error()}</h3>
                  </div>
                </div>
              </div>
            )}

            <div class="space-y-6">
              <div>
                <label for="username" class="block text-sm font-medium text-gray-11">
                  Username <span class="text-red-500">*</span>
                </label>
                <div class="mt-1">
                  <Input
                    type="text"
                    name="username"
                    id="username"
                    required
                    value={username()}
                    onInput={(e) => setUsername(e.currentTarget.value)}
                  />
                </div>
              </div>

              <div>
                <label for="password" class="block text-sm font-medium text-gray-11">
                  Password <span class="text-red-500">*</span>
                </label>
                <div class="mt-1">
                  <Input
                    type="password"
                    name="password"
                    id="password"
                    required
                    value={password()}
                    onInput={(e) => setPassword(e.currentTarget.value)}
                  />
                </div>
              </div>

              <div>
                <label for="confirmPassword" class="block text-sm font-medium text-gray-11">
                  Confirm Password <span class="text-red-500">*</span>
                </label>
                <div class="mt-1">
                  <Input
                    type="password"
                    name="confirmPassword"
                    id="confirmPassword"
                    required
                    value={confirmPassword()}
                    onInput={(e) => setConfirmPassword(e.currentTarget.value)}
                  />
                </div>
              </div>

              <div class="flex justify-end">
                <button
                  type="submit"
                  disabled={loading()}
                  class="inline-flex justify-center rounded-lg border border-transparent bg-blue-600 py-2 px-4 text-sm font-medium text-white shadow-sm hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 disabled:opacity-50"
                >
                  {loading() ? 'Creating...' : 'Create Login'}
                </button>
              </div>
            </div>
          </form>
        </div>
      </div>
    </div>
  );
};

export default CreateLogin;
