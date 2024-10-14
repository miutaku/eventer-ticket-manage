import PostalMime from 'postal-mime';

export default {
  async email(message, env, ctx) {
    try {
      // Parse the incoming email message using PostalMime
      const email = await PostalMime.parse(message.raw);

      console.log('Subject', email.subject);
      console.log('HTML', email.html);
      console.log('Text', email.text);

      // Check if the email has a 'from' field and if the specific sender is included
      if (email.from && email.from.value && email.from.value.some(sender => sender.address === 'lt-mail@l-tike.com')) {
        // Define the payload containing the email body
        const payload = {
          subject: email.subject,
          body: email.text || email.html || 'No content',
        };

        // Define the target endpoint where you want to send the payload
        const endpoint = 'https://apidebug.mtakumi-0925.workers.dev';

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
        console.log('The email is not from the specified sender.');
      }
    } catch (error) {
      console.error('An error occurred:', error);
    }
  },
};