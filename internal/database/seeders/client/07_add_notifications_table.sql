CREATE TABLE IF NOT EXISTS public.notifications (
     id SERIAL PRIMARY KEY,
     title VARCHAR(255) NOT NULL,
     message TEXT NOT NULL,
     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
     updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO notifications (title, message) VALUES
   ('New Feature Released', 'We just launched dark mode. Try it out now!'),
   ('Weekly Summary', 'Your weekly activity summary is now available.'),
   ('Security Alert', 'Unusual login detected from a new device.'),
   ('Post Approved', 'Your latest post has been approved and is now live.'),
   ('New Follower', 'Someone just followed your author profile.'),
   ('Comment Received', 'You have a new comment on your post.'),
   ('Maintenance Scheduled', 'System maintenance is scheduled for Sunday at 2AM UTC.'),
   ('Subscription Expiring', 'Your subscription will expire in 3 days. Renew now to avoid interruption.'),
   ('Welcome!', 'Thanks for joining our platform. Start creating your first post today.'),
   ('Beta Access Granted', 'You have been granted access to the beta features. Give us your feedback!');
