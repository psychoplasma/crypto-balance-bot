'use client';

import { useAuth } from '@/context/AuthContext';
import { useRouter } from 'next/navigation';
import { useState } from 'react';
import { login_ } from '@/actions/actions'
import './LoginForm.css';

export default function LoginForm() {
  const [error, setError] = useState('');
  const router = useRouter();
  const { login: authLogin } = useAuth();

  const submitAction = async (formData: FormData) => {
    setError('');

    try {
      const { userId, token } = await login_(
        formData.get('email') as string,
        formData.get('password') as string, //FIXME: pass password hash instead)
      );
      const res = await authLogin(userId, token);

      router.push(`/user?userId=${userId}`);
    } catch (error) {
      setError('An error occurred during login');
    }
  };

  return (
    <form action={submitAction}>
      <div className="form-group">
        <label htmlFor="email">Email:</label>
        <input type="email"  name="email" required/>
      </div>
      <div className="form-group">
        <label htmlFor="password">Password:</label>
        <input type="password" name="password" required/>
      </div>
      {error && <p className="error">{error}</p>}
      <button type="submit" className="submit-btn">Login</button>
    </form>
  );
};
