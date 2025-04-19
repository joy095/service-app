import type { LayoutLoad } from './$types';
import { redirect } from '@sveltejs/kit';

interface ParentData {
	user: any;
}

export const load: LayoutLoad = async ({ parent }) => {
	const { user } = (await parent()) as ParentData;
	if (!user) throw redirect(302, '/login');
	return {};
};
