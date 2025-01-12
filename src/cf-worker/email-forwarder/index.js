import PostalMime from 'postal-mime';

export default {
  async email(message, env, ctx) {
    try {
      // Parse the incoming email message using PostalMime
      const email = await PostalMime.parse(message.raw);

      console.log('Subject', email.subject);
      console.log('HTML', email.html);
      console.log('Text', email.text);

      // Check if the email body contains the specified text.
      const bodyContent = email.text || email.html || '';
      if (bodyContent) {
        // Define the payload containing the email body
        const payload = {
          subject: email.subject,
          body: bodyContent || 'No content',
        };

        // Retrieve the target endpoint from environment variables
        const endpoint = env.EMAIL_FORWARDER;

        // Send the POST request with the email content
        const response = await fetch(endpoint, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify(payload),
        });

        console.log('Notification sent with status:', response.status);
      } else {
        console.log('The email does not contain the specified text.');
      }
    } catch (error) {
      console.error('An error occurred:', error);
    }
  },
};
