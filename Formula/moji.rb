class Moji < Formula
  desc "Terminal art toolkit - kaomoji, ASCII banners, filters, QR codes, and more"
  homepage "https://github.com/ddmoney420/moji"
  version "1.0.1"
  license "MIT"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/ddmoney420/moji/releases/download/v#{version}/moji_darwin_arm64.tar.gz"
      sha256 "9f77efcd3fa132ae59c9db43625c53e6da0fc280a3ad2591946bae1f859e7fb4"
    else
      url "https://github.com/ddmoney420/moji/releases/download/v#{version}/moji_darwin_amd64.tar.gz"
      sha256 "f3e4428ed51c4c36b4da520219e5c87f26009cfc768d423da83fce4ec3f3bcc1"
    end
  end

  on_linux do
    if Hardware::CPU.arm?
      url "https://github.com/ddmoney420/moji/releases/download/v#{version}/moji_linux_arm64.tar.gz"
      sha256 "a34c4c2c8bd6aba8e2245bbcbd37e4dd7eabae903bf14619b8ba8b6b1d276149"
    else
      url "https://github.com/ddmoney420/moji/releases/download/v#{version}/moji_linux_amd64.tar.gz"
      sha256 "3931153c18786eac82fbb5cc03f73c90917e7570b406477f74e4d51d50d8148f"
    end
  end

  def install
    bin.install "moji"
  end

  test do
    assert_match "moji", shell_output("#{bin}/moji --version")
    assert_match "shrug", shell_output("#{bin}/moji list | head -20")
  end
end
