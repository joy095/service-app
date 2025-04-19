import { redirect } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ cookies }) => {
	const token = cookies.get('access_token');
	const user = cookies.get('user'); // optional

	if (token) {
		// User is already logged in, redirect to profile
		throw redirect(302, '/profile');
	}
};
