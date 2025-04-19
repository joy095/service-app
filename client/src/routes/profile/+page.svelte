<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { createQuery, createMutation } from '@tanstack/svelte-query';
	import { logout as apiLogout } from '$lib/api';
	import { API } from '$lib/share';

	onMount(() => {
		const token = localStorage.getItem('access_token');
		const user = JSON.parse(localStorage.getItem('user') || '{}');

		if (!token || !user?.username) {
			goto('/login'); // 🔒 redirect to login if not logged in
		}
	});

	// User profile query using TanStack Query
	const userQuery = createQuery({
		queryKey: ['userProfile'],
		queryFn: async () => {
			const accessToken = localStorage.getItem('access_token');
			const refreshToken = localStorage.getItem('refresh_token');
			const userFromStorage = JSON.parse(localStorage.getItem('user') || '{}');
			const username = userFromStorage?.username;

			if (!accessToken || !username) {
				throw new Error('Not logged in');
			}

			const fetchUserProfile = async (token: string) => {
				const res = await fetch(`${API}/v1/auth/user/${username}`, {
					method: 'GET',
					headers: {
						'Content-Type': 'application/json',
						Authorization: `Bearer ${token}`
					}
				});
				return res;
			};

			let res = await fetchUserProfile(accessToken);

			// If token expired or unauthorized
			if (res.status === 401 && refreshToken) {
				const refreshRes = await fetch(`${API}/v1/auth/refresh-token`, {
					method: 'POST',
					headers: {
						'Content-Type': 'application/json',
						Refresh_token: refreshToken
					}
				});

				const refreshData = await refreshRes.json();

				if (refreshRes.ok && refreshData.access_token && refreshData.refresh_token) {
					localStorage.setItem('access_token', refreshData.access_token);
					localStorage.setItem('refresh_token', refreshData.refresh_token);

					// Retry original request with new token
					res = await fetchUserProfile(refreshData.access_token);
				} else {
					// If refresh failed, logout
					localStorage.removeItem('access_token');
					localStorage.removeItem('refresh_token');
					localStorage.removeItem('user');
					goto('/login');
					throw new Error('Session expired. Please login again.');
				}
			}

			const data = await res.json();

			if (!res.ok) {
				throw new Error(data.message || 'Failed to fetch user profile');
			}

			return data.user;
		}
	});

	// Logout mutation
	const logoutMutation = createMutation({
		mutationFn: async () => {
			const accessToken = localStorage.getItem('access_token');
			const user = JSON.parse(localStorage.getItem('user') || '{}');

			if (!accessToken || !user?.id) {
				throw new Error('Already logged out');
			}

			return apiLogout(user.id, accessToken);
		},
		onSuccess: () => {
			// Clear local storage
			localStorage.removeItem('access_token');
			localStorage.removeItem('refresh_token');
			localStorage.removeItem('user');

			goto('/login'); // Redirect to login page
		},
		onError: (error) => {
			alert(`Failed to logout: ${error.message}`);
		}
	});

	function handleLogout() {
		$logoutMutation.mutate();
	}
</script>

{#if $userQuery.isPending}
	<p class="mt-10 text-center">Loading profile...</p>
{:else if $userQuery.isError}
	<p class="mt-10 text-center text-red-600">{$userQuery.error.message}</p>
{:else if $userQuery.data}
	<div class="mx-auto mt-10 max-w-md rounded border p-6 shadow">
		<h2 class="mb-4 text-xl font-bold">User Profile</h2>
		<p><strong>Username:</strong> {$userQuery.data.username}</p>
		<p><strong>Email:</strong> {$userQuery.data.email}</p>
		<p><strong>ID:</strong> {$userQuery.data.id}</p>
	</div>

	<button
		class="mt-4 rounded bg-red-600 px-4 py-2 text-white"
		on:click={handleLogout}
		disabled={$logoutMutation.isPending}
	>
		{$logoutMutation.isPending ? 'Logging out...' : 'Logout'}
	</button>
{/if}
