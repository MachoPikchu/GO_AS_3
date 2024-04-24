async function signup() {
    const email = document.getElementById('signup-email').value;
    const password = document.getElementById('signup-password').value;

    try {
        const response = await fetch('http://localhost:3000/signup', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ email, password })
        });

        if (!response.ok) {
            throw new Error('Signup failed');
        }

        document.getElementById('message').innerText = 'Signup successful!';
    } catch (error) {
        console.error(error);
        document.getElementById('message').innerText = 'Signup failed';
    }
}

async function login() {
    const email = document.getElementById('login-email').value;
    const password = document.getElementById('login-password').value;

    try {
        const response = await fetch('http://localhost:3000/login', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ email, password })
        });

        if (!response.ok) {
            throw new Error('Login failed');
        }

        const data = await response.json();
        document.getElementById('message').innerText = 'Login successful!';
        // Handle successful login, e.g., redirect to dashboard
    } catch (error) {
        console.error(error);
        document.getElementById('message').innerText = 'Login failed';
    }
}
