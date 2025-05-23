import SocialButton from "./social-button";

interface GoogleButtonProps {
  text?: string;
  className?: string;
  onSuccess?: () => void;
}

const GoogleButton = ({
  text = "Continue with Google",
  className,
  onSuccess,
}: GoogleButtonProps) => {
  const handleGoogleLogin = async () => {
    // Implementation for Google login goes here
    // This is a placeholder for the actual implementation
    console.log("Google login clicked");

    if (onSuccess) {
      onSuccess();
    }
  };

  return (
    <SocialButton
      icon={
        <img src="/google-black.svg" alt="Google logo" className="w-5 h-5" />
      }
      onClick={handleGoogleLogin}
      className={className}
    >
      {text}
    </SocialButton>
  );
};

export default GoogleButton;
