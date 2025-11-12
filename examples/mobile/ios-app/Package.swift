// swift-tools-version: 5.9
import PackageDescription

let package = Package(
    name: "AionMCPDemo",
    platforms: [
        .iOS(.v16)
    ],
    products: [
        .library(
            name: "AionMCPDemo",
            targets: ["AionMCPDemo"]
        )
    ],
    dependencies: [
        .package(url: "https://github.com/Alamofire/Alamofire.git", from: "5.8.0")
    ],
    targets: [
        .target(
            name: "AionMCPDemo",
            dependencies: ["Alamofire"]
        ),
        .testTarget(
            name: "AionMCPDemoTests",
            dependencies: ["AionMCPDemo"]
        )
    ]
)
