<script lang="ts">
	import { goto } from '$app/navigation';
	import { createMutation } from '@tanstack/svelte-query';
	import { apiFetch } from '$lib/api';

	let username = '';
	let email = '';
	let password = '';
	let confirmPassword = '';

	let hasSubmitted = false;

	let errors = {
		username: '',
		email: '',
		password: '',
		confirmPassword: ''
	};

	const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
	const usernameRegex = /^[a-zA-Z0-9_]{3,}$/;
	const passwordRegex = /^(?=.*[a-z])(?=.*[A-Z])(?=.*[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]).{6,}$/;

	let showPassword = false;
	let showConfirmPassword = false;
	let generalError = '';

	$: isFormValid =
		username &&
		email &&
		password &&
		confirmPassword &&
		!errors.username &&
		!errors.email &&
		!errors.password &&
		!errors.confirmPassword;

	function validateField(field: string, value: string) {
		switch (field) {
			case 'username':
				if (!value) errors.username = 'Username is required';
				else if (!usernameRegex.test(value))
					errors.username = 'Must be at least 3 characters, letters/numbers/_';
				else errors.username = '';
				break;

			case 'email':
				if (!value) errors.email = 'Email is required';
				else if (!emailRegex.test(value)) errors.email = 'Invalid email format';
				else errors.email = '';
				break;

			case 'password':
				if (!value) errors.password = 'Password is required';
				else if (!passwordRegex.test(value))
					errors.password =
						'Must contain uppercase, lowercase, special character, and be 6+ characters';
				else errors.password = '';
				validateField('confirmPassword', confirmPassword);
				break;

			case 'confirmPassword':
				if (!value) errors.confirmPassword = 'Confirm your password';
				else if (value !== password) errors.confirmPassword = 'Passwords do not match';
				else errors.confirmPassword = '';
				break;
		}
	}

	function validateAll() {
		validateField('username', username);
		validateField('email', email);
		validateField('password', password);
		validateField('confirmPassword', confirmPassword);
	}

	// Registration mutation with TanStack Query
	const registerMutation = createMutation({
		mutationFn: async (userData: { username: string; email: string; password: string }) => {
			const API = import.meta.env.VITE_API_URL;
			if (!API) {
				throw new Error('API URL is not defined');
			}

			const res = await apiFetch(`${API}/v1/auth/register`, {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify(userData)
			});

			const data = await res.json();

			if (!res.ok) {
				throw { data, status: res.status };
			}

			return data;
		},
		onSuccess: (data) => {
			const accessToken = data.tokens.access_token;
			const refreshToken = data.tokens.refresh_token;

			document.cookie = `access_token=${accessToken}; path=/; Secure; SameSite=Strict`;
			document.cookie = `refresh_token=${refreshToken}; path=/; Secure; SameSite=Strict`;

			localStorage.setItem('access_token', accessToken);
			localStorage.setItem('refresh_token', refreshToken);
			localStorage.setItem('user', JSON.stringify(data.user));

			// Modify fetch to include auth headers
			window.fetch = (
				(originalFetch) =>
				(input, init = {}) => {
					const headers = new Headers(init.headers || {});
					const token = localStorage.getItem('access_token');
					const refresh = localStorage.getItem('refresh_token');
					if (token) headers.set('Authorization', `Bearer ${token}`);
					if (refresh) headers.set('Refresh_token', refresh);

					return originalFetch(input, {
						...init,
						headers
					});
				}
			)(window.fetch);

			goto('/profile');
		},
		onError: (error: any) => {
			console.error('Registration error:', error);

			if (error.data) {
				if (error.data.errors) {
					errors = error.data.errors;
				}
				generalError = error.data.error || 'Something went wrong. Please try again.';
			} else {
				generalError = error.message || 'An unexpected error occurred';
			}
		}
	});

	async function handleSubmit(event: Event) {
		event.preventDefault();
		hasSubmitted = true;
		generalError = ''; // Clear previous error

		validateAll();

		if (Object.values(errors).some((e) => e)) return;

		// Submit the form using TanStack mutation
		$registerMutation.mutate({ username, email, password });
	}
</script>

<form on:submit={handleSubmit} class="login_form container mx-auto my-20">
	<h3>Register your account</h3>

	<!-- Username -->
	<div class="input_box">
		<label for="username">Username</label>
		<input
			id="username"
			bind:value={username}
			placeholder="Enter username"
			on:input={() => {
				if (hasSubmitted) validateField('username', username);
			}}
		/>
		{#if hasSubmitted && errors.username}
			<p class="text-red-500">{errors.username}</p>
		{/if}
	</div>

	<!-- Email -->
	<div class="input_box">
		<label for="email">Email</label>
		<input
			type="email"
			id="email"
			bind:value={email}
			placeholder="Enter email address"
			on:input={() => {
				if (hasSubmitted) validateField('email', email);
			}}
		/>
		{#if hasSubmitted && errors.email}
			<p class="text-red-500">{errors.email}</p>
		{/if}
	</div>

	<!-- Password -->
	<div class="input_box">
		<div class="password_title">
			<label for="password">Password</label>
		</div>
		<input
			id="password"
			type={showPassword ? 'text' : 'password'}
			bind:value={password}
			placeholder="Enter your password"
			on:input={() => {
				if (hasSubmitted) validateField('password', password);
			}}
		/>

		<span
			class="absolute bottom-0 right-5 top-0 my-auto h-6 w-6 cursor-pointer"
			on:click={() => (showPassword = !showPassword)}
		>
			{#if showPassword}
				<img src="/icon/eye-open.svg" alt="icon" loading="lazy" />
			{:else}
				<img src="/icon/eye-close.svg" alt="icon" loading="lazy" />
			{/if}
		</span>

		{#if hasSubmitted && errors.password}
			<p class="text-red-500">{errors.password}</p>
		{/if}
	</div>

	<!-- Confirm Password -->
	<div class="input_box relative">
		<div class="password_title">
			<label for="confirmPassword">Confirm Password</label>
		</div>
		<input
			type={showConfirmPassword ? 'text' : 'password'}
			id="confirmPassword"
			bind:value={confirmPassword}
			placeholder="Confirm your password"
			on:input={() => {
				if (hasSubmitted) validateField('confirmPassword', confirmPassword);
			}}
		/>

		<span
			class="absolute bottom-0 right-5 top-0 my-auto h-6 w-6 cursor-pointer"
			on:click={() => (showConfirmPassword = !showConfirmPassword)}
		>
			{#if showConfirmPassword}
				<img src="/icon/eye-open.svg" alt="icon" loading="lazy" />
			{:else}
				<img src="/icon/eye-close.svg" alt="icon" loading="lazy" />
			{/if}
		</span>

		{#if hasSubmitted && errors.confirmPassword}
			<p class="text-red-500">{errors.confirmPassword}</p>
		{/if}
	</div>

	<!-- Submit -->
	<button type="submit" disabled={!isFormValid || $registerMutation.isPending}>
		{$registerMutation.isPending ? 'Registering...' : 'Register'}
	</button>

	<p class="sign_up">Don't have an account? <a href="/login">Sign in</a></p>

	{#if generalError}
		<p class="mb-4 text-red-500">{generalError}</p>
	{/if}
</form>

<style>
	.login_form {
		width: 100%;
		max-width: 435px;
		background: #fff;
		border-radius: 6px;
		padding: 41px 30px;
		box-shadow: 0 10px 20px rgba(0, 0, 0, 0.15);
	}
	.login_form h3 {
		font-size: 20px;
		text-align: center;
	}

	form .input_box label {
		display: block;
		font-weight: 500;
		margin-bottom: 0.1rem;
	}
	/* Input field styling */
	form .input_box input {
		width: 100%;
		height: 3rem;
		border: 1px solid #dadaf2;
		border-radius: 5px;
		outline: none;
		background: #f8f8fb;
		font-size: 1rem;
		padding: 0px 1rem;
		margin-bottom: 1rem;
		transition: 0.2s ease;
	}
	form .input_box input:focus {
		border-color: #626cd6;
	}
	form .input_box .password_title {
		display: flex;
		justify-content: space-between;
		text-align: center;
	}
	form .input_box {
		position: relative;
	}
	a {
		text-decoration: none;
		color: #626cd6;
		font-weight: 500;
	}
	a:hover {
		text-decoration: underline;
	}
	/* Login button styling */
	form button {
		width: 100%;
		height: 3rem;
		border-radius: 5px;
		border: none;
		outline: none;
		background: #626cd6;
		color: #fff;
		font-size: 18px;
		font-weight: 500;
		text-transform: uppercase;
		cursor: pointer;
		margin-bottom: 28px;
		transition: 0.3s ease;
	}
	form button:hover {
		background: #4954d0;
	}

	form button:disabled {
		background: rgba(73, 84, 208, 0.7);
		cursor: not-allowed;
	}
</style>
