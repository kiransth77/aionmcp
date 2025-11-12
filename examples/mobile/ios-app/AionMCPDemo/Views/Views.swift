import SwiftUI

struct ToolsView: View {
    var body: some View {
        NavigationView {
            VStack(spacing: 20) {
                Text("Available Tools")
                    .font(.title)
                
                Text("Configure server URL in Settings to load tools")
                    .foregroundColor(.secondary)
                    .multilineTextAlignment(.center)
                    .padding()
                
                Spacer()
            }
            .padding()
            .navigationTitle("Tools")
        }
    }
}

struct StatsView: View {
    var body: some View {
        NavigationView {
            VStack(spacing: 20) {
                Text("Statistics")
                    .font(.title)
                
                Text("Server statistics will appear here")
                    .foregroundColor(.secondary)
                    .multilineTextAlignment(.center)
                    .padding()
                
                Spacer()
            }
            .padding()
            .navigationTitle("Statistics")
        }
    }
}

struct SettingsView: View {
    @State private var serverURL = ""
    @State private var apiKey = ""
    
    var body: some View {
        NavigationView {
            Form {
                Section("Server Configuration") {
                    TextField("Server URL", text: $serverURL)
                        .textInputAutocapitalization(.never)
                        .autocorrectionDisabled()
                        .keyboardType(.URL)
                    
                    SecureField("API Key (Optional)", text: $apiKey)
                }
                
                Section {
                    Button("Save") {
                        // Save settings
                    }
                    .frame(maxWidth: .infinity)
                }
            }
            .navigationTitle("Settings")
        }
    }
}

#Preview("Tools") {
    ToolsView()
}

#Preview("Stats") {
    StatsView()
}

#Preview("Settings") {
    SettingsView()
}
