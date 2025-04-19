<script lang="ts">
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';
	import { createMutation } from '@tanstack/svelte-query';
	import { API } from '$lib/share';

	onMount(() => {
		const token = localStorage.getItem('access_token');
		const user = JSON.parse(localStorage.getItem('user') || '{}');

		if (token && user?.username) {
			goto('/profile');
		}
	});

	let username = '';
	let password = '';
	let error: string | null = null;
	let showPassword = false;

	// Using TanStack Query mutation for login
	const mutation = createMutation({
		mutationFn: async (credentials: { username: string; password: string }) => {
			const res = await fetch(`${API}/v1/auth/login`, {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json'
				},
				body: JSON.stringify(credentials)
			});

			if (!res.ok) {
				const err = await res.json();
				throw new Error(err.message || 'Login failed');
			}

			return res.json();
		},
		onSuccess: (data) => {
			localStorage.setItem('access_token', data.tokens.access_token);
			localStorage.setItem('refresh_token', data.tokens.refresh_token);
			localStorage.setItem('user', JSON.stringify(data.user));

			// Navigate to a protected route
			goto('/profile');
		},
		onError: (err: Error) => {
			error = err.message;
		}
	});

	function login() {
		error = null;
		$mutation.mutate({ username, password });
	}
</script>

<div class="login-container">
	<h1>Login</h1>

	{#if error}
		<div class="error-message">{error}</div>
	{/if}

	<form on:submit|preventDefault={login}>
		<div class="form-group">
			<label for="username">Username</label>
			<input type="text" id="username" bind:value={username} required />
		</div>

		<div class="form-group">
			<label for="password">Password</label>
			<div class="password-field">
				<input
					type={showPassword ? 'text' : 'password'}
					id="password"
					bind:value={password}
					required
				/>
				<button type="button" on:click={() => (showPassword = !showPassword)}>
					{showPassword ? 'Hide' : 'Show'}
				</button>
			</div>
		</div>

		<button type="submit" disabled={$mutation.isPending}>
			{$mutation.isPending ? 'Logging in...' : 'Login'}
		</button>
	</form>
</div>

<style>
	.login-container {
		max-width: 400px;
		margin: 0 auto;
		padding: 2rem;
		border-radius: 8px;
		box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
	}

	h1 {
		text-align: center;
		margin-bottom: 1.5rem;
	}

	.form-group {
		margin-bottom: 1rem;
	}

	label {
		display: block;
		margin-bottom: 0.5rem;
		font-weight: 500;
	}

	input {
		width: 100%;
		padding: 0.5rem;
		border: 1px solid #ccc;
		border-radius: 4px;
	}

	.password-field {
		display: flex;
	}

	.password-field input {
		flex: 1;
		border-radius: 4px 0 0 4px;
	}

	.password-field button {
		padding: 0.5rem;
		background: #f1f1f1;
		border: 1px solid #ccc;
		border-left: none;
		border-radius: 0 4px 4px 0;
		cursor: pointer;
	}

	button[type='submit'] {
		width: 100%;
		padding: 0.75rem;
		background: #4a4e69;
		color: white;
		border: none;
		border-radius: 4px;
		cursor: pointer;
		font-size: 1rem;
		margin-top: 1rem;
	}

	button[type='submit']:hover {
		background: #3d405b;
	}

	button[type='submit']:disabled {
		background: #9a9ca9;
		cursor: not-allowed;
	}

	.error-message {
		background: #ffebee;
		color: #c62828;
		padding: 0.75rem;
		border-radius: 4px;
		margin-bottom: 1rem;
	}
</style>
