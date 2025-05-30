import { Component, createSignal } from 'solid-js';
import { userApi } from '../api/user';
import { A } from '@solidjs/router';

const FindUsername: Component = () => {
  const [email, setEmail] = createSignal('');
  const [error, setError] = createSignal<string | null>(null);
  const [success, setSuccess] = createSignal(false);
  const [loading, setLoading] = createSignal(false);

  const handleSubmit = async (e: Event) => {
    e.preventDefault();
    setError(null);
    setSuccess(false);
    setLoading(true);

    try {
      await userApi.findUsername(email());
      setSuccess(true);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to find username');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div class="min-h-screen bg-gray-1 flex flex-col justify-center py-12 sm:px-6 lg:px-8">
      <div class="sm:mx-auto sm:w-full sm:max-w-md">
        <h2 class="mt-6 text-center text-3xl font-extrabold text-gray-12">
          Find your username
        </h2>
        <p class="mt-2 text-center text-sm text-gray-11">
          Enter your email address and we'll send you your username
        </p>
      </div>

      <div class="mt-8 sm:mx-auto sm:w-full sm:max-w-md">
        <div class="bg-white py-8 px-4 shadow-lg rounded-lg sm:px-10">
          <form class="space-y-6" onSubmit={handleSubmit}>
            <div>
              <label
                for="email"
                class="block text-sm font-medium text-gray-11"
              >
                Email address
              </label>
              <div class="mt-1">
                <input
                  id="email"
                  name="email"
                  type="email"
                  autocomplete="email"
                  required
                  value={email()}
                  onInput={(e) => setEmail(e.currentTarget.value)}
                  class="appearance-none block w-full px-3 py-2 border border-gray-7 rounded-lg shadow-sm placeholder:text-gray-8 focus:outline-none focus:ring-2 focus:ring-primary-7 focus:border-primary-7"
                />
              </div>
            </div>

            {error() && (
              <div class="rounded-lg bg-red-2 p-4">
                <div class="flex">
                  <div class="ml-3">
                    <h3 class="text-sm font-medium text-red-11">{error()}</h3>
                  </div>
                </div>
              </div>
            )}

            {success() && (
              <div class="rounded-lg bg-green-2 p-4">
                <div class="flex">
                  <div class="ml-3">
                    <h3 class="text-sm font-medium text-green-11">
                      If an account exists with that email, we will send the username to it.
                    </h3>
                  </div>
                </div>
              </div>
            )}

            <div>
              <button
                type="submit"
                disabled={loading()}
                class="w-full flex justify-center py-2 px-4 border border-transparent rounded-lg shadow-sm text-sm font-medium text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50 disabled:cursor-not-allowed"
              >
                {loading() ? 'Sending...' : 'Send username'}
              </button>
            </div>

            <div class="text-center">
              <A
                href="/login"
                class="text-sm text-blue-600 hover:text-blue-500"
              >
                Back to login
              </A>
            </div>
          </form>
        </div>
      </div>
    </div>
  );
};

export default FindUsername;
