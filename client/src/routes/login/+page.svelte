<script lang="ts">
	let username = '';
	let password = '';

	let hasSubmitted = false;

	let errors = {
		username: '',
		password: ''
	};

	let showPassword = false;

	$: isLoginFormValid = username && password && !errors.username && !errors.password;

	async function handleSubmit(event: Event) {
		event.preventDefault();
		hasSubmitted = true;

		if (Object.values(errors).some((e) => e)) return;

		const res = await fetch('api/v1/auth/login', {
			method: 'POST',
			headers: { 'Content-Type': 'application/json' },
			body: JSON.stringify({ username, password })
		});

		const data = await res.json();

		if (!res.ok) {
			errors = data.errors ?? {};
		} else {
			alert('Login successful!');
		}
	}
</script>

<form on:submit={handleSubmit} class="login_form container mx-auto my-20">
	<h3>Log in with</h3>

	<!-- Username -->
	<div class="input_box">
		<label for="username">Username</label>
		<input id="username" bind:value={username} placeholder="Enter username" />
		{#if hasSubmitted && errors.username}
			<p class="text-red-500">{errors.username}</p>
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
		/>

		<span
			class="absolute top-0 right-5 bottom-0 my-auto h-6 w-6 cursor-pointer"
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

	<!-- Submit -->
	<button type="submit" disabled={!isLoginFormValid}>Log In</button>
	<p class="sign_up">Don't have an account? <a href="/register">Sign up</a></p>
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
