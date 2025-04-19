import { API } from './share';

export async function apiFetch(path: string, options: RequestInit = {}) {
	let token = localStorage.getItem('access_token');
	const refresh = localStorage.getItem('refresh_token');

	const headers = new Headers(options.headers || {});
	if (token) headers.set('Authorization', `Bearer ${token}`);
	if (refresh) headers.set('Refresh_token', refresh);

	let res = await fetch(path, {
		...options,
		headers
	});

	if (res.status === 401) {
		// try refreshing
		const refreshRes = await fetch(`${API}/v1/auth/refresh-token`, {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json',
				Refresh_token: refresh || ''
			}
		});
		const data = await refreshRes.json();
		if (refreshRes.ok) {
			token = data.access_token;
			if (token) {
				localStorage.setItem('access_token', token);
				headers.set('Authorization', `Bearer ${token}`);
			}
			res = await fetch(path, {
				...options,
				headers
			});
		}
	}

	return res;
}

export async function logout(user_id: string, access_token: string) {
	try {
		const res = await fetch(`${API}/v1/auth/logout`, {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json',
				Authorization: `Bearer ${access_token}`
			},
			body: JSON.stringify({ user_id })
		});

		if (!res.ok) {
			const error = await res.json();
			throw new Error(error.message || 'Logout failed');
		}
	} catch (err) {
		console.error('Logout error:', err);
		throw err;
	}
}
