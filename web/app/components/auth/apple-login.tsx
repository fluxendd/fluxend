import SocialButton from "./social-button";

interface AppleButtonProps {
  text?: string;
  className?: string;
  onSuccess?: () => void;
}

const AppleButton = ({
  text = "Continue with Apple",
  className,
  onSuccess,
}: AppleButtonProps) => {
  const handleAppleLogin = async () => {
    console.log("Apple login clicked");

    if (onSuccess) {
      onSuccess();
    }
  };

  return (
    <SocialButton
      icon={<img src="/apple.svg" alt="Apple logo" className="w-5 h-5 mr-2" />}
      onClick={handleAppleLogin}
      className={className}
    >
      {text}
    </SocialButton>
  );
};

export default AppleButton;
