<?php

namespace Tests\Helper;

use Illuminate\Console\Command;
use Illuminate\Console\OutputStyle;
use Mockery as m;
use Mockery\MockInterface;
use Symfony\Component\Console\Input\InputInterface;

/**
 * Trait ConsoleHelper
 *
 * @package Tests\Helper
 */
trait ConsoleHelper
{
    /**
     * @param Command        $command
     * @param InputInterface $input
     * @param OutputStyle    $output
     *
     * @return Command
     */
    private function setInputAndOutput(Command $command, InputInterface $input, OutputStyle $output): Command
    {
        $command->setInput($input);
        $command->setOutput($output);

        return $command;
    }

    /**
     * @return InputInterface|MockInterface
     */
    private function createInput(): InputInterface
    {
        return m::spy(InputInterface::class);
    }

    /**
     * @param InputInterface|MockInterface $input
     * @param string|null                  $argument
     * @param string|null                  $key
     *
     * @return $this
     */
    private function mockInputGetArgument(MockInterface $input, ?string $argument, ?string $key): self
    {
        $input
            ->shouldReceive('getArgument')
            ->with($key)
            ->andReturn($argument);

        return $this;
    }

    /**
     * @return OutputStyle|MockInterface
     */
    private function createOutputStyle(): OutputStyle
    {
        return m::spy(OutputStyle::class);
    }

    /**
     * @param OutputStyle|MockInterface $outputStyle
     * @param string                    $line
     *
     * @param int|null                  $verbosity
     *
     * @return $this
     */
    private function assertOutputStyleWriteln(MockInterface $outputStyle, string $line, int $verbosity): self
    {
        $outputStyle
            ->shouldHaveReceived('writeln')
            ->with($line, $verbosity)
            ->once();

        return $this;
    }
}
