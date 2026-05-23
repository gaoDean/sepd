document.addEventListener('DOMContentLoaded', () => {
  const form = document.getElementById('contact-form');
  const submitBtn = document.getElementById('submit-btn');
  const statusDiv = document.getElementById('form-status');

  if (!form) return;

  form.addEventListener('submit', async (e) => {
    e.preventDefault();

    // Clear previous status messages
    statusDiv.style.display = 'none';
    statusDiv.textContent = '';
    statusDiv.style.backgroundColor = '';
    statusDiv.style.color = '';
    statusDiv.style.border = '';

    // Disable submit button and show loading text
    submitBtn.disabled = true;
    const originalBtnText = submitBtn.textContent;
    submitBtn.textContent = 'Sending...';

    // Retrieve and trim values
    const name = document.getElementById('name').value.trim();
    const email = document.getElementById('email').value.trim();
    const message = document.getElementById('message').value.trim();

    // Basic client-side validation
    if (!name || !email || !message) {
      showStatus('All fields are required.', 'error');
      resetSubmitButton(originalBtnText);
      return;
    }

    const payload = { name, email, message };

    // Resolve API endpoint: since the Go backend runs on port 8080 on the server,
    // we direct the requests to port 8080 of the active hostname.
    const protocol = window.location.protocol === 'https:' ? 'https:' : 'http:';
    const hostname = window.location.hostname || '127.0.0.1';
    const apiEndpoint = `${protocol}//${hostname}:8080/api/contact`;

    try {
      const response = await fetch(apiEndpoint, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(payload)
      });

      const contentType = response.headers.get('content-type');
      let data = {};
      if (contentType && contentType.includes('application/json')) {
        data = await response.json();
      }

      if (response.ok) {
        showStatus('Your message has been sent successfully! We will get back to you soon.', 'success');
        form.reset();
      } else {
        throw new Error(data.error || `Failed to send message (Status: ${response.status})`);
      }
    } catch (err) {
      showStatus(err.message || 'An error occurred. Please try again later.', 'error');
    } finally {
      resetSubmitButton(originalBtnText);
    }
  });

  function showStatus(text, type) {
    statusDiv.style.display = 'block';
    statusDiv.textContent = text;
    statusDiv.style.padding = '0.75rem';
    statusDiv.style.marginTop = '1rem';
    
    if (type === 'success') {
      // Styled matching positive notification tones (greenish)
      statusDiv.style.backgroundColor = 'rgba(56, 142, 60, 0.1)';
      statusDiv.style.color = '#2e7d32';
      statusDiv.style.border = '1px solid #2e7d32';
    } else {
      // Styled matching negative notification tones (reddish)
      statusDiv.style.backgroundColor = 'rgba(198, 40, 40, 0.1)';
      statusDiv.style.color = '#c62828';
      statusDiv.style.border = '1px solid #c62828';
    }
  }

  function resetSubmitButton(text) {
    submitBtn.disabled = false;
    submitBtn.textContent = text;
  }
});
