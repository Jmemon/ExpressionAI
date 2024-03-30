import AVFoundation

class TextToSpeechConverter {
    let speechSynthesizer = AVSpeechSynthesizer()

    func speak(text: String) {
        let speechUtterance = AVSpeechUtterance(string: text)
        speechUtterance.voice = AVSpeechSynthesisVoice(language: "en-US")
        speechSynthesizer.speak(speechUtterance)
    }
}
