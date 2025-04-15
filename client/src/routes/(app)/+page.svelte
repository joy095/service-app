<nav class="navbar">
	<div class="navbar-left">
		<a href="/" class="logo"><img src="images/logo.png" alt="logo" /></a>
		<div class="search-box">
			<img src="images/search.png" alt="search" />
			<input type="text" placeholder="Search for anything" />
		</div>
	</div>
	<div class="navbar-center">
		<ul>
			<li>
				<a href="#" class="active-link"
					><img src="images/home.png" alt="home" /> <span>Home</span></a
				>
			</li>
			<li>
				<a href="#"><img src="images/network.png" alt="network" /> <span>My Network</span></a>
			</li>
			<li><a href="#"><img src="images/jobs.png" alt="jobs" /> <span>Jobs</span></a></li>
			<li>
				<a href="#"><img src="images/message.png" alt="message" /> <span>Messaging</span></a>
			</li>
			<li>
				<a href="#"
					><img src="images/notification.png" alt="notification" />
					<span>Notifications</span></a
				>
			</li>
		</ul>
	</div>
	<div class="navbar-right">
		<div class="online">
			<img src="images/user-1.png" class="nav-profile-img" />
		</div>
	</div>
	<!----Dropdown menu-->
	<div class="profile-menu-wrap" id="profileMenu">
		<div class="profile-menu">
			<div class="user-info">
				<img src="images/user-1.png" />
				<div>
					<h3>John Doe</h3>
					<a href="#">See your profile</a>
				</div>
			</div>
			<hr />
			<a href="#" class="profile-menu-link">
				<img src="images/feedback.png" />
				<p>Give Feedback</p>
				<span>></span>
			</a>
			<a href="#" class="profile-menu-link">
				<img src="images/setting.png" />
				<p>Settings & Privacy</p>
				<span>></span>
			</a>
			<a href="#" class="profile-menu-link">
				<img src="images/help.png" />
				<p>Help & Support</p>
				<span>></span>
			</a>
			<a href="#" class="profile-menu-link">
				<img src="images/display.png" />
				<p>Display & Accessibility</p>
				<span>></span>
			</a>
			<a href="#" class="profile-menu-link">
				<img src="images/logout.png" />
				<p>Logout</p>
				<span>></span>
			</a>
		</div>
	</div>
</nav>

<div class="container">
	<div class="left-sidebar">
		<div class="sidebar-profile-box">
			<img src="images/cover-pic.png" width="100%" />
			<div class="sidebar-profile-info">
				<img src="images/user-1.png" />
				<h1>John Doe</h1>
				<h3>Data Analyst at IBM</h3>
				<ul>
					<li>Your profile views <span>24</span></li>
					<li>Your post views <span>128</span></li>
					<li>Your Connections<span>108</span></li>
				</ul>
			</div>
			<div class="sidebar-profile-link">
				<a href="#"><img src="images/items.png" />My Items</a>
				<a href="#"><img src="images/premium.png" />Try Premium</a>
			</div>
		</div>

		<div class="sidebar-activity" id="sidebarActivity">
			<h3>RECENT</h3>
			<a href="#"><img src="images/recent.png" />Data Analysis</a>
			<a href="#"><img src="images/recent.png" />UI UX Design</a>
			<a href="#"><img src="images/recent.png" />Web Development</a>
			<a href="#"><img src="images/recent.png" />Object Oriented Programming</a>
			<a href="#"><img src="images/recent.png" />Operating Systems</a>
			<a href="#"><img src="images/recent.png" />Platform technologies</a>
			<h3>GROUPS</h3>
			<a href="#"><img src="images/group.png" />Data Analyst group</a>
			<a href="#"><img src="images/group.png" />Learn NumPy</a>
			<a href="#"><img src="images/group.png" />Machine Learning group</a>
			<a href="#"><img src="images/group.png" />Data Science Aspirants</a>
			<h3>HASHTAG</h3>
			<a href="#"><img src="images/hashtag.png" />dataanalyst</a>
			<a href="#"><img src="images/hashtag.png" />numpy</a>
			<a href="#"><img src="images/hashtag.png" />machinelearning</a>
			<a href="#"><img src="images/hashtag.png" />datascience</a>
			<div class="discover-more-link">
				<a href="#">Discover More</a>
			</div>
		</div>
		<p id="showMoreLink">Show more <b>+</b></p>
	</div>
	<div class="main-content">
		<div class="create-post">
			<div class="create-post-input">
				<img src="images/user-1.png" />
				<textarea rows="2" placeholder="Write Something"></textarea>
			</div>
			<div class="create-post-links">
				<li><img src="images/photo.png" />Photo</li>
				<li><img src="images/video.png" />Video</li>
				<li><img src="images/event.png" />Event</li>
				<li>Post</li>
			</div>
		</div>
		<div class="sort-by">
			<hr />
			<p>Sort by : <span>top <img src="images/down-arrow.png" /></span></p>
		</div>
		<div class="post">
			<div class="post-author">
				<img src="images/user-3.png" />
				<div>
					<h1>Bejamin Leo</h1>
					<small>Founder and CEO at Giva | Angel Investor</small>
					<small>2 hours ago </small>
				</div>
			</div>
			<p>
				The sucess of every website depends on Search engine optimisation and digital marketing
				strategy. If you are not in the first page of all major search engines then you are ahead
				among your competitors.
			</p>
			<img src="images/post-image-1.png" width="100%" />

			<div class="post-stats">
				<div>
					<img src="images/thumbsup.png" />
					<img src="images/love.png" />
					<img src="images/clap.png" />
					<span class="liked-users">Adam Doe and 89 others</span>
				</div>
				<div>
					<span>22 comments &middot; 40 shares</span>
				</div>
			</div>
			<div class="post-activity">
				<div>
					<img src="images/user-1.png" class="post-activity-user-icon" />
					<img src="images/down-arrow.png" class="post-activity-arrow-icon" />
				</div>
				<div class="post-activity-link">
					<img src="images/like.png" />
					<span>Like</span>
				</div>
				<div class="post-activity-link">
					<img src="images/comment.png" />
					<span>Comment</span>
				</div>
				<div class="post-activity-link">
					<img src="images/share.png" />
					<span>Share</span>
				</div>
				<div class="post-activity-link">
					<img src="images/send.png" />
					<span>Send</span>
				</div>
			</div>
		</div>

		<div class="post">
			<div class="post-author">
				<img src="images/user-4.png" />
				<div>
					<h1>Claire Smith</h1>
					<small>SDE at Swiggy | Solopreneur</small>
					<small>2 hours ago </small>
				</div>
			</div>
			<p>
				The sucess of every website depends on Search engine optimisation and digital marketing
				strategy. If you are not in the first page of all major search engines then you are ahead
				among your competitors.
			</p>
			<img src="images/post-image-2.png" width="100%" />

			<div class="post-stats">
				<div>
					<img src="images/thumbsup.png" />
					<img src="images/love.png" />
					<img src="images/clap.png" />
					<span class="liked-users">Adam Doe and 89 others</span>
				</div>
				<div>
					<span>22 comments &middot; 40 shares</span>
				</div>
			</div>
			<div class="post-activity">
				<div>
					<img src="images/user-1.png" class="post-activity-user-icon" />
					<img src="images/down-arrow.png" class="post-activity-arrow-icon" />
				</div>
				<div class="post-activity-link">
					<img src="images/like.png" />
					<span>Like</span>
				</div>
				<div class="post-activity-link">
					<img src="images/comment.png" />
					<span>Comment</span>
				</div>
				<div class="post-activity-link">
					<img src="images/share.png" />
					<span>Share</span>
				</div>
				<div class="post-activity-link">
					<img src="images/send.png" />
					<span>Send</span>
				</div>
			</div>
		</div>

		<div class="post">
			<div class="post-author">
				<img src="images/user-3.png" />
				<div>
					<h1>Bejamin Leo</h1>
					<small>Founder and CEO at Giva | Angel Investor</small>
					<small>2 hours ago </small>
				</div>
			</div>
			<p>
				The sucess of every website depends on Search engine optimisation and digital marketing
				strategy. If you are not in the first page of all major search engines then you are ahead
				among your competitors.
			</p>
			<img src="images/post-image-3.png" width="100%" />

			<div class="post-stats">
				<div>
					<img src="images/thumbsup.png" />
					<img src="images/love.png" />
					<img src="images/clap.png" />
					<span class="liked-users">Adam Doe and 89 others</span>
				</div>
				<div>
					<span>22 comments &middot; 40 shares</span>
				</div>
			</div>
			<div class="post-activity">
				<div>
					<img src="images/user-1.png" class="post-activity-user-icon" />
					<img src="images/down-arrow.png" class="post-activity-arrow-icon" />
				</div>
				<div class="post-activity-link">
					<img src="images/like.png" />
					<span>Like</span>
				</div>
				<div class="post-activity-link">
					<img src="images/comment.png" />
					<span>Comment</span>
				</div>
				<div class="post-activity-link">
					<img src="images/share.png" />
					<span>Share</span>
				</div>
				<div class="post-activity-link">
					<img src="images/send.png" />
					<span>Send</span>
				</div>
			</div>
		</div>

		<div class="post">
			<div class="post-author">
				<img src="images/user-3.png" />
				<div>
					<h1>Bejamin Leo</h1>
					<small>Founder and CEO at Giva | Angel Investor</small>
					<small>2 hours ago </small>
				</div>
			</div>
			<p>
				The sucess of every website depends on Search engine optimisation and digital marketing
				strategy. If you are not in the first page of all major search engines then you are ahead
				among your competitors.
			</p>
			<img src="images/post-image-4.png" width="100%" />

			<div class="post-stats">
				<div>
					<img src="images/thumbsup.png" />
					<img src="images/love.png" />
					<img src="images/clap.png" />
					<span class="liked-users">Adam Doe and 89 others</span>
				</div>
				<div>
					<span>22 comments &middot; 40 shares</span>
				</div>
			</div>
			<div class="post-activity">
				<div>
					<img src="images/user-1.png" class="post-activity-user-icon" />
					<img src="images/down-arrow.png" class="post-activity-arrow-icon" />
				</div>
				<div class="post-activity-link">
					<img src="images/like.png" />
					<span>Like</span>
				</div>
				<div class="post-activity-link">
					<img src="images/comment.png" />
					<span>Comment</span>
				</div>
				<div class="post-activity-link">
					<img src="images/share.png" />
					<span>Share</span>
				</div>
				<div class="post-activity-link">
					<img src="images/send.png" />
					<span>Send</span>
				</div>
			</div>
		</div>
	</div>
	<div class="right-sidebar">
		<div class="sidebar-news">
			<img src="images/more.png" class="info-icon" />
			<h3>Trending News</h3>

			<a href="#">High Demand for Skilled Employees</a>
			<span>1d ago &middot; 10,934 readers</span>

			<a href="#">Inflation in Canada Affects the Workforce</a>
			<span>2d ago &middot; 7,043 readers</span>

			<a href="#">Mass Recruiters fire Employees</a>
			<span>4d ago &middot; 17,789 readers</span>

			<a href="#">Crypto predicted to Boom this year</a>
			<span>9d ago &middot; 2, 436 readers</span>

			<a href="#" class="read-more-link">Read More</a>
		</div>

		<div class="sidebar-ad">
			<small>Ad &middot; &middot; &middot;</small>
			<p>Master Web Development</p>
			<div>
				<img src="images/user-1.png" />
				<img src="images/mi-logo.png" />
			</div>
			<b>Brand and Demand in Xiaomi</b>
			<a href="#" class="ad-link">Learn More</a>
		</div>

		<div class="sidebar-useful-links">
			<a href="#">About</a>
			<a href="#">Accessibility</a>
			<a href="#">Help Center</a>
			<a href="#">Privacy Policy</a>
			<a href="#">Advertising</a>
			<a href="#">Get the App</a>
			<a href="#">More</a>

			<div class="copyright-msg">
				<img src="images/logo.png" />
				<p>Linkedin &#169; 2022. All Rights Reserved</p>
			</div>
		</div>
	</div>
</div>

<style>
	* {
		margin: 0;
		padding: 0;
		font-family: 'Poppins', sans-serif;
		box-sizing: border-box;
	}
	body {
		background-color: #f0f2f5;
		color: #5f5f5f;
	}
	a {
		text-decoration: none;
		color: #5f5f5f;
	}
	.navbar {
		display: flex;
		align-items: center;
		justify-content: space-between;
		background: #fff;
		padding: 6px 6%;
		position: sticky;
		top: 0;
		z-index: 100;
		box-shadow: 0 5px 10px rgba(0, 0, 0, 0.1);
	}
	.logo img {
		width: 30px;
		margin-right: 15px;
		display: block;
	}
	.navbar-center ul li {
		display: inline-block;
		list-style: none;
	}

	.navbar-center ul li a {
		display: flex;
		align-items: center;
		font-size: 14px;
		margin: 5px 8px;
		padding-right: 5px;
		position: relative;
	}
	.navbar-center ul li a img {
		width: 30px;
	}
	.navbar-center ul li a::after {
		content: '';
		width: 0;
		height: 2px;
		background: #045be6;
		position: absolute;
		bottom: -10px;
		transition: width 0.3s;
	}
	.navbar-center ul li a:hover::after,
	.navbar-center ul li a.active-link::after {
		width: 100%;
	}
	.nav-profile-img {
		width: 40px;
		border-radius: 50%;
		display: block;
		cursor: pointer;
		position: relative;
	}

	.online {
		position: relative;
	}
	.online::after {
		content: '';
		width: 7px;
		height: 7px;
		border: 2px solid #fff;
		border-radius: 50%;
		right: 0;
		top: 0;
		background: #41db51;
		position: absolute;
	}
	.search-box {
		background: #f0f2f5;
		width: 250px;
		border-radius: 20px;
		display: flex;
		align-items: center;
		padding: 0px 15px;
	}
	.navbar-left {
		display: flex;
		align-items: center;
	}
	.search-box img {
		width: 14px;
	}
	.search-box input {
		width: 100%;
		background: transparent;
		padding: 8px;
		outline: none;
		border: none;
	}

	/*   */

	.container {
		display: flex;
		padding: 20px 6%;
		flex-wrap: wrap;
		justify-content: space-between;
	}

	/* Left side Bar*/
	.left-sidebar {
		flex-basis: 25%;
		align-self: flex-start;
		position: sticky;
		top: 73px;
	}

	.sidebar-profile-box {
		background: #fff;
	}
	.sidebar-profile-info {
		padding: 0 25px;
	}
	.sidebar-profile-info img {
		width: 90px;
		border-radius: 50%;
		background: #fff;
		padding: 4px;
		margin-top: -45px;
	}
	.sidebar-profile-info h1 {
		font-size: 20px;
		font-weight: 600;
		color: #222;
	}
	.sidebar-profile-info h3 {
		font-size: 14px;
		font-weight: 500;
		color: #777;
	}
	.sidebar-profile-info ul {
		list-style: none;
		margin: 20px 0;
	}

	.sidebar-profile-info ul li {
		width: 100%;
		margin: 5px 0;
		font-size: 13px;
	}

	.sidebar-profile-info ul li span {
		float: right;
		color: #045be6;
	}
	.sidebar-profile-link {
		display: flex;
		align-items: center;
		border-top: 1px solid #ccc;
	}
	.sidebar-profile-link a {
		flex-basis: 50%;
		display: flex;
		align-items: center;
		justify-content: center;
		padding: 15px 0;
		font-size: 13px;
		border-left: 1px solid #ccc;
	}
	.sidebar-profile-link a:first-child {
		border-left: 0;
	}
	.sidebar-profile-link a img {
		width: 20px;
		margin-right: 10px;
	}
	.sidebar-activity {
		background: #fff;
		padding: 5px 25px;
		margin: 12px 0;
	}
	.sidebar-activity h3 {
		color: #777;
		font-size: 14px;
		font-weight: 500;
		margin: 20px 0 10px;
	}
	.sidebar-activity a {
		display: flex;
		align-items: center;
		font-size: 12px;
		font-weight: 500;
		color: #888;
		margin: 3px 0;
	}
	.sidebar-activity a img {
		width: 20px;
		margin-right: 10px;
	}
	.discover-more-link {
		border-top: 1px solid #ccc;
		text-align: center;
		margin-top: 20px;
		margin-left: -25px;
		margin-right: -25px;
	}
	.discover-more-link a {
		color: #045be6;
		display: inline-block;
		margin: 10px 0;
	}
	/*Right Sidebar */

	.right-sidebar {
		flex-basis: 25%;
		align-self: flex-start;
		position: sticky;
		top: 73px;
	}
	.sidebar-news {
		background: #fff;
		padding: 10px 25px;
	}
	.info-icon {
		width: 15px;
		float: right;
		margin-top: 15px;
	}
	.sidebar-news h3 {
		font-size: 18 px;
		font-weight: 600;
		color: #333;
		margin: 10px 0 30px;
	}
	.sidebar-news a {
		display: block;
		font-size: 12px;
		font-weight: 600;
		margin-top: 10px;
		margin-bottom: -5px;
	}
	.sidebar-news span {
		font-size: 12px;
	}

	.sidebar-news .read-more-link {
		color: #045be6;
		font-weight: 500;
		margin: 20px 0 10px;
	}
	.sidebar-ad {
		background: #ffff;
		padding: 15px 25px;
		text-align: center;
		margin: 12px 0;
		font-size: 12px;
	}

	.sidebar-ad img {
		width: 60px;
		border-radius: 50%;
		margin: 4px;
	}
	.sidebar-ad small {
		float: right;
		font-weight: 500;
	}
	.sidebar-ad p {
		margin-top: 30px;
		margin-bottom: 10px;
	}

	.sidebar-ad b {
		display: block;
		font-weight: 500;
		margin-top: 10px;
	}
	.ad-link {
		display: inline-block;
		border: 1px solid #045be6;
		border-radius: 30px;
		padding: 5px 15px;
		color: #045be6;
		font-weight: 500;
		margin: 20px auto 10px;
	}

	.sidebar-useful-links {
		padding: 15px 25px;
		text-align: center;
	}
	.sidebar-useful-links a {
		display: inline-block;
		font-size: 14px;
		margin: 4px 2px;
	}
	.copyright-msg {
		display: flex;
		align-items: center;
		justify-content: center;
		font-size: 12px;
		margin-top: 10px;
		font-weight: 500;
	}
	.copyright-msg img {
		width: 18px;
		margin-right: 4px;
	}

	/*Main Content*/

	.main-content {
		flex-basis: 47%;
	}
	.create-post {
		background: #fff;
	}
	.create-post-input {
		padding: 20px 25px 10px;
		display: flex;
		align-items: flex-start;
	}
	.create-post-input img {
		width: 35px;
		border-radius: 50%;
		margin-right: 10px;
	}
	.create-post-input textarea {
		width: 100%;
		border: 0;
		outline: 0;
		resize: none;
		background: transparent;
		margin-top: 8px;
	}
	::placeholder {
		font-weight: 500;
	}
	.create-post-links {
		display: flex;
		align-items: flex-start;
	}
	.create-post-links li {
		list-style: none;
		border-top: 1px solid #ccc;
		border-right: 1px solid #ccc;
		flex-basis: 25%;
		height: 40px;
		font-size: 13px;
		display: flex;
		align-items: center;
		justify-content: center;
		cursor: pointer;
	}
	.create-post-links li img {
		width: 15px;
		margin-right: 5px;
	}
	.create-post-links li:last-child {
		background: #045be6;
		color: #fff;
		border-top: 0;
		border-right: 0;
	}
	.sort-by {
		display: flex;
		align-items: center;
		margin: 10px 0;
	}
	.sort-by hr {
		flex: 1;
		border: 0;
		height: 1px;
		background: #ccc;
	}
	.sort-by p {
		font-size: 13px;
		padding-left: 5px;
	}
	.sort-by p span {
		font-weight: 600;
		cursor: pointer;
	}
	.sort-by p span img {
		width: 12px;
		margin-left: 3px;
	}
	.post {
		background: #fff;
		padding: 20px 25px 5px;
		margin: 5px 0 15px;
	}
	.post-author {
		display: flex;
		align-items: flex-start;
		margin-bottom: 20px;
	}
	.post-author img {
		width: 35px;
		border-radius: 100px;
		margin-top: 5px;
		margin-right: 5px;
	}
	.post-author h1 {
		font-size: 18px;
		font-weight: 600;
		color: #333;
	}
	.post-author small {
		display: block;
		margin-bottom: -2px;
	}
	.post p {
		font-size: 14px;
		margin-bottom: 15px;
	}
	.post-stats {
		display: flex;
		flex-wrap: wrap;
		align-items: center;
		justify-content: space-between;
		font-size: 12px;
		border-bottom: 1px solid #ccc;
		padding-bottom: 6px;
	}
	.post-stats div {
		display: flex;
		align-items: center;
	}

	.post-stats img {
		width: 15px;
		margin-right: -5px;
	}
	.liked-users {
		margin-left: 10px;
	}
	.post-activity {
		display: flex;
		align-items: center;
		justify-content: space-between;
		padding: 10px;
		font-size: 14px;
		font-weight: 500;
	}
	.post-activity div {
		display: flex;
		align-items: center;
	}
	.post-activity-user-icon {
		width: 22px;
		border-radius: 50%;
	}
	.post-activity-arrow-icon {
		width: 12px;
		margin-left: 5px;
	}
	.post-activity-link img {
		width: 18px;
		margin-right: 8px;
	}

	.profile-menu-wrap {
		position: absolute;
		top: 100%;
		right: 5%;
		width: 320px;
		max-height: 0;
		overflow: hidden;
		transition: max-height 0.5s;
	}
	.profile-menu-wrap.open-menu {
		max-height: 400px;
	}
	.profile-menu {
		background: #222;
		color: #fff;
		padding: 20px;
		margin: 10px;
	}
	.user-info {
		display: flex;
		align-items: center;
	}
	.user-info img {
		width: 50px;
		border-radius: 50%;
		margin-right: 15px;
	}
	.user-info h3 {
		font-weight: 500;
		margin-bottom: 0 px;
	}
	.user-info a {
		color: #c9dbf8;
		font-size: 13px;
	}
	.profile-menu hr {
		border: 0;
		height: 1px;
		width: 100%;
		background: #fff;
		margin: 15px 0 10px;
	}
	.profile-menu-link {
		display: flex;
		align-items: center;
		color: #fff;
		margin: 12px 0;
		font-size: 14px;
	}
	.profile-menu-link p {
		width: 100%;
	}
	.profile-menu-link img {
		width: 35px;
		border-radius: 50%;
		margin-right: 15px;
	}
	.profile-menu-link span {
		font-size: 10px;
	}

	#showMoreLink {
		display: none;
	}

	/* MEDIA QUERIES */
	@media only screen and (max-width: 600px) {
		.search-box {
			background: transparent;
			width: auto;
			padding: 0;
		}
		.search-box input {
			width: 0;
			padding: 0;
		}
		.navbar-center ul li a span {
			display: none;
		}
		.navbar-center ul li a {
			padding-right: 0;
			margin: 5px;
		}
		.nav-profile-img {
			width: 30px;
		}
		.container {
			padding: 15px 3%;
		}
		.left-sidebar,
		.right-sidebar {
			flex-basis: 100%;
			position: relative;
			top: unset;
		}
		.main-content {
			flex-basis: 100%;
		}
		.post-stats div {
			margin: 2px 0;
		}
		.post-activity-link span {
			display: none;
		}
		.post-activity-link img {
			margin-right: 0;
		}
		#showMoreLink {
			display: block;
			width: fit-content;
			margin: 10px auto;
			font-size: 14px;
			cursor: pointer;
		}
		.sidebar-activity {
			display: none;
		}
		.sidebar-activity.open-activity {
			display: block;
		}
	}
</style>
