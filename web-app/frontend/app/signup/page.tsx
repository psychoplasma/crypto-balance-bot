
import Layout from '@/components/Layout';
import SignUpForm from '@/components/SignUpForm';
import PublicRoute from '@/components/PublicRoute';
import './signup.css';

const Signup = () => {
  return (
    <PublicRoute>
      <Layout>
        <div className={"signup-container"}>
          <div className={"signup-card"}>
            <h1>Sign Up</h1>
            <SignUpForm />
            <p className={"login-link"}>
              Already have an account? <a href="/login">Log in</a>
            </p>
          </div>
        </div>
      </Layout>
    </PublicRoute>
  );
};

export default Signup;
